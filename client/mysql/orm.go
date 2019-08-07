package mysql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/caijinlin/golib/helper"
	"github.com/caijinlin/golib/log"
	"github.com/jinzhu/gorm"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

func NewGorm(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dsn)
	if err == nil {
		db.LogMode(true)
		db.SetLogger(&GormLogger{})
		db.BlockGlobalUpdate(true) //avoid Missing WHERE clause while deleting/updating
	}
	return db, err
}

func NewSql(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	return db, err
}

type GormLogger struct {
}

// Print format & print log
func (logger GormLogger) Print(values ...interface{}) {

	// 	s.print("sql", fileWithLineNum(), NowFunc().Sub(t), sql, vars, s.RowsAffected)

	if len(values) > 1 {
		level := values[0]
		position := fmt.Sprintf("%v", values[1])
		if level == "sql" {
			log.Info(map[string]interface{}{
				"action":   "mysql_call",
				"sql":      FormatedSql(values...),
				"cost":     helper.FormatDurationToMs(values[2].(time.Duration)),
				"errmsg":   fmt.Sprintf("%v", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned "),
				"position": position,
			})
		} else {
			log.Warning(map[string]interface{}{
				"action":   "mysql_call",
				"errmsg":   fmt.Sprintf("%s", values[2:]...),
				"position": position,
			})
		}
	}
}

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

func FormatedSql(values ...interface{}) (sql string) {
	var formattedValues []string
	for _, value := range values[4].([]interface{}) {
		indirectValue := reflect.Indirect(reflect.ValueOf(value))
		if indirectValue.IsValid() {
			value = indirectValue.Interface()
			if t, ok := value.(time.Time); ok {
				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
			} else if b, ok := value.([]byte); ok {
				if str := string(b); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			} else if r, ok := value.(driver.Valuer); ok {
				if value, err := r.Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			} else {
				switch value.(type) {
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
					formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
				default:
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}
		} else {
			formattedValues = append(formattedValues, "NULL")
		}
	}

	if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
		sql = values[3].(string)
		for index, value := range formattedValues {
			placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
			sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
		}
	} else {
		formattedValuesLength := len(formattedValues)
		for index, value := range sqlRegexp.Split(values[3].(string), -1) {
			sql += value
			if index < formattedValuesLength {
				sql += formattedValues[index]
			}
		}
	}

	return sql
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

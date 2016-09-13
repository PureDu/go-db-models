package models

import (
	"strings"
)

type ParseTable struct {
	Name        string // 数据表名，如：ad_plan，对应结构体名为：AdPlanTable，对应文件名为：gen_ad_plan.go
	PackageName string // 生成的程序的包名
	PrimaryType string // 主键的类型，如：uint32, sql.NullString等
	Fields      []ParseField
}

type ParseField struct {
	Name string // 字段名，如：plan_id。对应到结构体的属性名就是：PlanId
	Type string // 字段类型，对应golang中的类型，如：uint32, sql.NullString
}

// 解释数据表的结构体
func ParseTablesStruct(tables []Table, package_name string) (parse_tables []ParseTable) {
	for _, table := range tables {
		ptable := ParseTable{
			Name:        table.Name,
			PackageName: package_name,
		}
		ptable.Fields, ptable.PrimaryType = ParseFieldsStruct(table.Fields)

		// 生成代码文件
		GenFile(ptable)

		parse_tables = append(parse_tables, ptable)
	}

	return parse_tables
}

// 解释一个数据表的所有字段
func ParseFieldsStruct(fields []Field) (pfields []ParseField, primary_type string) {
	for _, f := range fields {
		pf := ParseField{
			Name: f.Name,
		}

		if isString(f.Type) {
			// 字符串
			if f.Null == "YES" {
				pf.Type = "sql.NullString"
			} else {
				pf.Type = "string"
			}
		} else if strings.Contains(f.Type, "int") {
			// 整型
			if f.Null == "YES" {
				pf.Type = "sql.NullInt64"
			} else {
				prefix := ""
				if strings.Contains(f.Type, "unsigned") {
					prefix = "u"
				}

				int_type := "int32"
				if strings.Contains(f.Type, "tinyint") {
					int_type = "int8"
				} else if strings.Contains(f.Type, "smallint") {
					int_type = "int16"
				} else if strings.Contains(f.Type, "bigint") {
					int_type = "int64"
				}
				pf.Type = prefix + int_type
			}
		} else if strings.Contains(f.Type, "float") {
			// 浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
			} else {
				pf.Type = "float"
			}
		} else if strings.Contains(f.Type, "double") || strings.Contains(f.Type, "decimal") {
			// 高精度浮点数
			if f.Null == "YES" {
				pf.Type = "sql.NullFloat64"
			} else {
				pf.Type = "float64"
			}
		} else if strings.Contains(f.Type, "year") {
			// 年份
			if f.Null == "YES" {
				pf.Type = "sql.NullInt64"
			} else {
				pf.Type = "uint16"
			}
		} else if strings.Contains(f.Type, "datetime") || strings.Contains(f.Type, "timestamp") || strings.Contains(f.Type, "date") {
			// 日期时间
			pf.Type = "time.Time"
		}

		pfields = append(pfields, pf)
	}

	return pfields, primary_type
}

// 判断是否是字符串
func isString(field_type string) bool {
	return strings.Contains(field_type, "text") || strings.Contains(field_type, "char") || strings.Contains(field_type, "binary") || strings.Contains(field_type, "blob")
}

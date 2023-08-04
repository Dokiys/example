package basic

type Table struct {
	TableName    *string // 表名称
	PK           *string // 主键
	TableComment *string // 表注释
}

func (t *Table) GetTableName() string {
	if t.TableName == nil {
		return ""
	}
	return *t.TableName
}
func (t *Table) GetPK() string {
	if t.PK == nil {
		return ""
	}
	return *t.PK
}
func (t *Table) GetTableComment() string {
	if t.TableComment == nil {
		return ""
	}
	return *t.TableComment
}

type Column struct {
	TableCatalog           *string // 表目录
	TableSchema            *string // 表模式
	TableName              *string // 表名称
	ColumnName             *string // 字段名称
	OrdinalPosition        *int    // 字段顺序
	ColumnDefault          *string // 字段默认值
	IsNullable             *string // 是否可空
	DataType               *string // 字段类型
	CharacterMaximumLength *int    // 字段最大长度
	CharacterOctetLength   *int    // 字段最大长度
	NumericPrecision       *int    // 数字精度
	NumericScale           *int    // 数字精度
	DatetimePrecision      *int    // 时间精度
	CharacterSetName       *string // 字符集
	CollationName          *string // 字符集排序规则
	ColumnType             *string // 字段类型
	ColumnKey              *string // 字段键
	Extra                  *string // 额外信息
	Privileges             *string // 权限
	ColumnComment          *string // 字段注释
	GenerationExpression   *string // 字段生成表达式
	SrsId                  *int    // 空间参考系
}

func (c *Column) GetTableCatalog() string {
	if c.TableCatalog == nil {
		return ""
	}
	return *c.TableCatalog
}
func (c *Column) GetTableSchema() string {
	if c.TableSchema == nil {
		return ""
	}
	return *c.TableSchema
}
func (c *Column) GetTableName() string {
	if c.TableName == nil {
		return ""
	}
	return *c.TableName
}
func (c *Column) GetColumnName() string {
	if c.ColumnName == nil {
		return ""
	}
	return *c.ColumnName
}
func (c *Column) GetOrdinalPosition() int {
	if c.OrdinalPosition == nil {
		return 0
	}
	return *c.OrdinalPosition
}
func (c *Column) GetColumnDefault() string {
	if c.ColumnDefault == nil {
		return ""
	}
	return *c.ColumnDefault
}
func (c *Column) GetIsNullable() string {
	if c.IsNullable == nil {
		return ""
	}
	return *c.IsNullable
}
func (c *Column) GetDataType() string {
	if c.DataType == nil {
		return ""
	}
	return *c.DataType
}
func (c *Column) GetCharacterMaximumLength() int {
	if c.CharacterMaximumLength == nil {
		return 0
	}
	return *c.CharacterMaximumLength
}
func (c *Column) GetCharacterOctetLength() int {
	if c.CharacterOctetLength == nil {
		return 0
	}
	return *c.CharacterOctetLength
}
func (c *Column) GetNumericPrecision() int {
	if c.NumericPrecision == nil {
		return 0
	}
	return *c.NumericPrecision
}
func (c *Column) GetNumericScale() int {
	if c.NumericScale == nil {
		return 0
	}
	return *c.NumericScale
}
func (c *Column) GetDatetimePrecision() int {
	if c.DatetimePrecision == nil {
		return 0
	}
	return *c.DatetimePrecision
}
func (c *Column) GetCharacterSetName() string {
	if c.CharacterSetName == nil {
		return ""
	}
	return *c.CharacterSetName
}
func (c *Column) GetCollationName() string {
	if c.CollationName == nil {
		return ""
	}
	return *c.CollationName
}
func (c *Column) GetColumnType() string {
	if c.ColumnType == nil {
		return ""
	}
	return *c.ColumnType
}
func (c *Column) GetColumnKey() string {
	if c.ColumnKey == nil {
		return ""
	}
	return *c.ColumnKey
}
func (c *Column) GetExtra() string {
	if c.Extra == nil {
		return ""
	}
	return *c.Extra
}
func (c *Column) GetPrivileges() string {
	if c.Privileges == nil {
		return ""
	}
	return *c.Privileges
}
func (c *Column) GetColumnComment() string {
	if c.ColumnComment == nil {
		return ""
	}
	return *c.ColumnComment
}
func (c *Column) GetGenerationExpression() string {
	if c.GenerationExpression == nil {
		return ""
	}
	return *c.GenerationExpression
}
func (c *Column) GetSrsId() int {
	if c.SrsId == nil {
		return 0
	}
	return *c.SrsId
}

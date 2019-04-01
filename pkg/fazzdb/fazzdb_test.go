package fazzdb

var builder = Builder{}
var param = NewParameter(Config{
	Limit:  15,
	Offset: 5,
	Lock:   LO_FOR_SHARE,
})

type AITest struct {
	*Model
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Number int    `db:"number"`
}
type UTest struct {
	*Model
	Id     string `db:"id"`
	Name   string `db:"name"`
	Number int    `db:"number"`
}
type PTest struct {
	*Model
	Id     string `db:"id"`
	Name   string `db:"name"`
	Number int    `db:"number"`
}
type TsTest struct {
	*Model
	Id int `db:"id"`
}
type SdTest struct {
	*Model
	Id int `db:"id"`
}
type CpTest struct {
	Model
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Number int    `db:"number"`
}

var ai = &AITest{
	Model: AutoIncrementModel(
		"auto_increment_tests",
		[]Column{
			Col("id"),
			Col("name"),
			Col("number"),
		},
		"id",
		false,
		false,
	),
	Id:     1,
	Name:   "ai_name_test",
	Number: 20,
}
var u = &UTest{
	Model: UuidModel(
		"uuid_tests",
		[]Column{
			Col("id"),
			Col("name"),
			Col("number"),
		},
		"id",
		false,
		false,
	),
	Id:     "1c8da90a-dd85-4afa-a162-e9a9bddc0e81",
	Name:   "u_name_test",
	Number: 15,
}
var p = &PTest{
	Model: PlainModel(
		"plain_tests",
		[]Column{
			Col("id"),
			Col("name"),
			Col("number"),
		},
		"id",
		false,
		false,
	),
	Id:     "3271040101920001",
	Name:   "p_name_test",
	Number: 23,
}
var ts = &TsTest{
	Model: AutoIncrementModel(
		"timestamp_tests",
		[]Column{
			Col("id"),
		},
		"id",
		true,
		false,
	),
	Id: 2,
}
var sd = &SdTest{
	Model: AutoIncrementModel(
		"soft_delete_tests",
		[]Column{
			Col("id"),
		},
		"id",
		false,
		true,
	),
	Id: 3,
}
var c = NewCpTest()

func (cm *CpTest) Get(key string) interface{} {
	return cm.Payload()[key]
}

func (cm *CpTest) Payload() map[string]interface{} {
	return cm.MapPayload(cm)
}

func (um *UTest) GeneratePK() {
	um.GenerateId(um)
}

func NewCpTest() *CpTest {
	return &CpTest{
		Model: *AutoIncrementModel(
			"complete_tests",
			[]Column{
				Col("id"),
				Col("name"),
				Col("number"),
			},
			"id",
			true,
			true,
		),
	}
}

func initDefaultCpTest() {
	c.Id = 4
	c.Name = "c_name_test"
	c.Number = 10
}

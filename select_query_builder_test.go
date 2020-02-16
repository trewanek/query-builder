package query_builder

import (
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID string `db:"user_id" table:"users"`
	Name   string `db:"name" table:"users"`
	Age    int    `db:"age" table:"users"`
	Sex    string `db:"sex" table:"users"`
}

func Test_SelectQueryBuilder_OnlyTable(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Build()

	expected := "SELECT users.* FROM users;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Model(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Model(User{}).
		Build()

	expected := "SELECT users.user_id, users.name, users.age, users.sex FROM users;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Column(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Column("name", "age", "sex").
		Build()
	expected := "SELECT users.name, users.age, users.sex FROM users;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_OrderBy(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		OrderBy("created", Asc).
		Build()
	expected := "SELECT users.* FROM users ORDER BY created ASC;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Table("users").
		OrderBy("created, user_id", Desc).
		Build()
	expected2 := "SELECT users.* FROM users ORDER BY created, user_id DESC;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_GroupBy(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		GroupBy("user_id").
		Build()

	expected := "SELECT users.* FROM users GROUP BY user_id;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Limit(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Limit().
		Build()
	expected := "SELECT users.* FROM users LIMIT ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Table("users").
		Placeholder(Named).
		Limit().
		Build()
	expected2 := "SELECT users.* FROM users LIMIT :limit;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Offset(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Offset().
		Build()
	expected := "SELECT users.* FROM users OFFSET ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Table("users").
		Placeholder(Named).
		Offset().
		Build()
	expected2 := "SELECT users.* FROM users OFFSET :offset;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_OrderBy_GroupBy_Limit_Offset(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		OrderBy("created", Asc).
		GroupBy("user_id").
		Limit().
		Offset().
		Build()
	expected := "SELECT users.* FROM users GROUP BY user_id ORDER BY created ASC LIMIT ? OFFSET ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		OrderBy("created", Desc).
		GroupBy("user_id").
		Limit().
		Offset().
		Build()
	expected2 := "SELECT users.* FROM users GROUP BY user_id ORDER BY created DESC LIMIT :limit OFFSET :offset;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Where(t *testing.T) {
	// ? bind
	q := NewSelectQueryBuilder().Table("users").
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()
	expected := "SELECT users.* FROM users WHERE name = ? AND age >= ? AND age <= ? AND sex != ? AND age < ? AND age > ?;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	// column name bind
	q2 := NewSelectQueryBuilder().Table("users").
		Placeholder(Named).
		Where("name", Equal).
		Where("age", GraterEqual).
		Where("age", LessEqual).
		Where("sex", Not).
		Where("age", LessThan).
		Where("age", GraterThan).
		Build()
	expected2 := "SELECT users.* FROM users WHERE name = :name AND age >= :age AND age <= :age AND sex != :sex AND age < :age AND age > :age;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}

	// custom name bind
	q3 := NewSelectQueryBuilder().Table("users").
		Placeholder(Named).
		Where("name", Equal).
		Where("age", GraterEqual, "age1").
		Where("age", LessEqual, "age2").
		Where("sex", Not, "sex1").
		Where("age", LessThan, "age3").
		Where("age", GraterThan, "age4").
		Build()
	expected3 := "SELECT users.* FROM users WHERE name = :name AND age >= :age1 AND age <= :age2 AND sex != :sex1 AND age < :age3 AND age > :age4;"
	if err := checkQuery(expected3, q3); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_WhereIn(t *testing.T) {
	q := NewSelectQueryBuilder().
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected := "SELECT users.* FROM users WHERE user_name = ? AND user_id IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		Where("user_name", Equal).
		WhereIn("user_id", 3).
		Build()

	expected2 := "SELECT users.* FROM users WHERE user_name = :user_name AND user_id IN (:user_id1, :user_id2, :user_id3);"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_WhereNotIn(t *testing.T) {
	q := NewSelectQueryBuilder().Table("users").
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected := "SELECT users.* FROM users WHERE user_name = ? AND user_id NOT IN (?, ?, ?);"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	q2 := NewSelectQueryBuilder().Table("users").
		Placeholder(Named).
		Where("user_name", Equal).
		WhereNotIn("user_id", 3).
		Build()
	expected2 := "SELECT users.* FROM users WHERE user_name = :user_name AND user_id NOT IN (:user_id1, :user_id2, :user_id3);"
	if q2 != expected2 {
		t.Logf("expected: %s, acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_Join(t *testing.T) {
	joinFields := []string{"user_id"}
	q := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		Join(LeftJoin, "tasks", joinFields, joinFields).
		Build()
	expected := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	joinFields2 := []string{"user_id"}
	joinFields3 := []string{"task_id"}
	q2 := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("users").
		Join(LeftJoin, "tasks", joinFields2, joinFields2).
		Join(LeftJoin, "subtasks", joinFields3, joinFields3, "tasks").
		Build()
	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.user_id LEFT JOIN subtasks ON tasks.task_id = subtasks.task_id;"
	if err := checkQuery(expected2, q2); err != nil {
		t.Log(err)
		t.Fail()
	}

	// TODO add other joins
}

func Test_SelectQueryBuilder_WhereMultiByStruct(t *testing.T) {
	type SearchMachinesParameter struct { //ex Tagged struct
		MachineNumber *int       `search:"machine_number" operator:"eq"`
		MachineName   *string    `search:"machine_name" operator:"eq"`
		BuyDateFrom   *time.Time `search:"buy_date" operator:"ge"`
		BuyDateTo     *time.Time `search:"buy_date" operator:"lt"`
		PriceFrom     *int       `search:"price" operator:"gt"`
		PriceTo       *int       `search:"price" operator:"le"`
		Owner         *string    `search:"owner" operator:"not"`
	}

	machineNumber := 150
	machineName := "machine1"
	price := 1000
	now := time.Now()
	owner := "owner1"

	searchParam := SearchMachinesParameter{
		MachineNumber: &machineNumber,
		MachineName:   &machineName,
		BuyDateFrom:   &now,
		BuyDateTo:     &now,
		PriceFrom:     &price,
		PriceTo:       &price,
		Owner:         &owner,
	}

	q := NewSelectQueryBuilder().
		Placeholder(Named).
		Table("machines").
		WhereMultiByStruct(searchParam).
		Build()

	expected := "SELECT machines.* FROM " +
		"machines " +
		"WHERE machine_number = :machine_number " +
		"AND machine_name = :machine_name " +
		"AND buy_date >= :buy_date_from " +
		"AND buy_date < :buy_date_to " +
		"AND price > :price_from " +
		"AND price <= :price_to " +
		"AND owner != :owner;"

	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_JoinMultipleFields(t *testing.T) {
	fields := []string{"user_id", "task_id"}
	q := NewSelectQueryBuilder().Table("users").
		Join(LeftJoin, "tasks", fields, fields).
		Build()
	expected := "SELECT users.* FROM users " +
		"LEFT JOIN tasks " +
		"ON users.user_id = tasks.user_id AND users.task_id = tasks.task_id;"
	if err := checkQuery(expected, q); err != nil {
		t.Log(err)
		t.Fail()
	}

	// JOIN 先と元でField名が異なる場合のJOIN
	originFields := []string{"user_id", "user_task_id"}
	targetFields := []string{"task_user_id", "task_id"}
	q2 := NewSelectQueryBuilder().Table("users").
		Join(LeftJoin, "tasks", originFields, targetFields).
		Build()
	expected2 := "SELECT users.* FROM users LEFT JOIN tasks ON users.user_id = tasks.task_user_id AND users.user_task_id = tasks.task_id;"

	if q2 != expected2 {
		t.Logf("expected: %s\n acctual: %s", expected2, q2)
		t.Fail()
	}
}

func Test_SelectQueryBuilder_IsImmutable(t *testing.T) {
	qb := NewSelectQueryBuilder().
		Table("users").
		Offset()

	copied := qb.
		Table("tasks")

	if reflect.DeepEqual(qb, copied) {
		t.Fail()
		t.Log(qb, copied, " are deepEqual true. object is not immutable.")
	}
}
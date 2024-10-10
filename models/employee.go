package models

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./employees.db")
	if err != nil {
		return err
	}
	DB = db
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS `employee` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `name` text,`phone` text,`dept` text)")
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

type Employee struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Dept  string `json:"dept"`
}

func GetEmployees(count int) ([]Employee, error) {

	rows, err := DB.Query("SELECT * from employee LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	employees := make([]Employee, 0)

	for rows.Next() {
		employee := Employee{}
		err = rows.Scan(&employee.Id, &employee.Name, &employee.Phone, &employee.Dept)

		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return employees, err
}

func GetEmployee(id string) (Employee, error) {

	stmt, err := DB.Prepare("SELECT * from employee WHERE id = ?")

	if err != nil {
		return Employee{}, err
	}

	employee := Employee{}

	sqlErr := stmt.QueryRow(id).Scan(&employee.Id, &employee.Name, &employee.Phone, &employee.Dept)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Employee{}, nil
		}
		return Employee{}, sqlErr
	}
	return employee, nil
}

func SaveOrUpdateEmployee(newEmployee Employee) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}
	if newEmployee.Id != 0 {
		stmt, err := tx.Prepare("UPDATE employee SET name = ?, phone = ?, dept = ? WHERE id = ?")

		if err != nil {
			return false, err
		}

		defer stmt.Close()

		_, err = stmt.Exec(newEmployee.Name, newEmployee.Phone, newEmployee.Dept, newEmployee.Id)
		if err != nil {
			return false, err
		}
	} else {
		stmt, err := tx.Prepare("INSERT INTO employee (name, phone, dept) VALUES (?, ?, ?)")

		if err != nil {
			return false, err
		}

		defer stmt.Close()

		_, err = stmt.Exec(newEmployee.Name, newEmployee.Phone, newEmployee.Dept)
		if err != nil {
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

func DeleteEmployee(id int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from employee where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

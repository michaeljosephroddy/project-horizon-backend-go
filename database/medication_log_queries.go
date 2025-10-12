package database

//nb need to fix this it should join the medication_log_item table
var medicationLogQuery = `select * from medication_log where user_id = ? and date(taken_at) between ? and ?;`

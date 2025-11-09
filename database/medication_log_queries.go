package database

// nb need to fix this it should join the medication_log_item table
var medicationLogQuery = `SELECT 
		ml.medication_log_id,
		ml.user_id,
		ml.taken_at,
		ml.notes AS log_notes,
		ml.created_at,
		ml.updated_at,
		JSON_ARRAYAGG(
			JSON_OBJECT(
				'medication_id', m.medication_id,
				'name', m.name,
				'dosage', mli.dosage
			)
		) AS medications
	FROM medication_log ml
	JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
	JOIN medication m ON m.medication_id = mli.medication_id
	WHERE ml.user_id = ?
		AND DATE(ml.taken_at) BETWEEN ? AND ?
	GROUP BY 
		ml.medication_log_id,
		ml.user_id,
		ml.taken_at,
		ml.notes,
		ml.created_at,
		ml.updated_at
	ORDER BY ml.taken_at DESC`

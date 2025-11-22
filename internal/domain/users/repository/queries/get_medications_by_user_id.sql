SELECT 
    um.user_medication_id,
    um.medication_id,
    um.dosage,
    um.start_date,
    um.note,
    m.name,
    m.description
FROM user_medication um
JOIN medication m ON um.medication_id = m.medication_id
WHERE um.user_id = ? AND um.stopped = 0
ORDER BY m.name

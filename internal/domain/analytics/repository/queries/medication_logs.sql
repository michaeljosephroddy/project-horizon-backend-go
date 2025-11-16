select ml.medication_log_id,
       ml.user_id,
       ml.taken_at,
       ml.notes as log_notes,
       ml.created_at,
       ml.updated_at,
       json_arrayagg(json_object('medication_id', m.medication_id, 'name', m.name, 'dosage', mli.dosage)) as medications
from medication_log ml
join medication_log_item mli on ml.medication_log_id = mli.medication_log_id
join medication m on m.medication_id = mli.medication_id
where ml.user_id = ?
    and date(ml.taken_at) between ? and ?
group by ml.medication_log_id,
         ml.user_id,
         ml.taken_at,
         ml.notes,
         ml.created_at,
         ml.updated_at
order by ml.taken_at desc

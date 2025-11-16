select distinct date(ml.taken_at) as log_date
from medication_log ml
join medication_log_item mli on ml.medication_log_id = mli.medication_log_id
where ml.user_id = ?
    and mli.medication_id = ?
    and date(ml.taken_at) between ? and ?
order by log_date

select m.medication_id,
       m.name,
       count(*) as total_doses,
       count(distinct date(ml.taken_at)) as days_active,
       datediff(?, ?) + 1 as total_days
from medication_log ml
join medication_log_item mli on ml.medication_log_id = mli.medication_log_id
join medication m on mli.medication_id = m.medication_id
where ml.user_id = ?
    and date(ml.taken_at) between ? and ?
group by m.medication_id,
         m.name
order by total_doses desc

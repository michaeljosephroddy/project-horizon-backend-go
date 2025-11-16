select count(distinct date(ml.taken_at)) as days_with_logs,
       datediff(?, ?) + 1 as total_days,
       avg(med_count) as avg_meds_per_log
from medication_log ml
left join
    (select medication_log_id,
            count(*) as med_count
     from medication_log_item
     group by medication_log_id) med_counts on ml.medication_log_id = med_counts.medication_log_id
left join medication_log_item mli on ml.medication_log_id = mli.medication_log_id
where ml.user_id = ?
    and date(ml.taken_at) between ? and ?

select time_format(sec_to_time(avg(time_to_sec(time(ml.taken_at)))), '%H:%i:%s') as avg_time,
       std(time_to_sec(time(ml.taken_at))) / 60 as std_dev_minutes,
       time_format(min(time(ml.taken_at)), '%H:%i:%s') as earliest_time,
       time_format(max(time(ml.taken_at)), '%H:%i:%s') as latest_time
from medication_log ml
join medication_log_item mli on ml.medication_log_id = mli.medication_log_id
where ml.user_id = ?
    and mli.medication_id = ?
    and date(ml.taken_at) between ? and ?

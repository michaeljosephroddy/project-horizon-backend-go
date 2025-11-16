select sl.sleep_log_id,
       sl.user_id,
       sl.hours_slept,
       sqt.name as sleep_quality_tag_name,
       sl.notes,
       sl.sleep_date,
       sl.created_at,
       sl.updated_at
from sleep_log sl
join sleep_quality_tag sqt on sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
where sl.user_id = ?
    and sl.sleep_date between ? and ?
    order  by sl.sleep_date;

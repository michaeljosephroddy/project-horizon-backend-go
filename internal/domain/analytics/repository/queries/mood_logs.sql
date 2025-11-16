select *
from mood_log
where user_id = ?
    and date(created_at) between ? and ?;

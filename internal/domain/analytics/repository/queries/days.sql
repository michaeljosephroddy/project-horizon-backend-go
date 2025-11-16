with mood_data as
    (select date(ml.created_at) as date,
            ml.created_at,
            ml.mood_log_id,
            ml.mood_rating,
            ml.note,
            group_concat(mt.name
                         order by mt.name separator ', ') as mood_tags,
            group_concat(mt.mood_tag_id
                         order by mt.mood_tag_id separator ',') as mood_tag_ids,
            sum(case
                    when mt.mood_category_id = ? then 1
                    else 0
                end) as entry_target_count,
            count(mt.mood_tag_id) as entry_total_count
     from mood_log ml
     join mood_log_mood_tag mlmt on ml.mood_log_id = mlmt.mood_log_id
     join mood_tag mt on mlmt.mood_tag_id = mt.mood_tag_id
     where ml.user_id = ?
         and date(ml.created_at) between ? and ?
     group by date(ml.created_at),
              ml.mood_log_id,
              ml.created_at,
              ml.mood_rating,
              ml.note),
     daily_stats as
    (select date, created_at,
                  mood_log_id,
                  mood_rating,
                  note,
                  mood_tags,
                  mood_tag_ids,
                  avg(mood_rating) over (partition by date) as daily_avg_rating,
                  sum(entry_target_count) over (partition by date) as daily_target_count,
                  sum(entry_total_count) over (partition by date) as daily_total_count,
                  (sum(entry_target_count) over (partition by date) * 100.0 / nullif(sum(entry_total_count) over (partition by date), 0)) as daily_target_percentage
     from mood_data)
select date, created_at,
             mood_log_id,
             mood_rating,
             note,
             mood_tags,
             mood_tag_ids,
             daily_avg_rating,
             daily_target_count,
             daily_total_count,
             coalesce(daily_target_percentage, 0) as daily_target_percentage
from daily_stats
where daily_avg_rating % s ?
    and coalesce(daily_target_percentage, 0) >= ?
order by date desc, created_at desc;

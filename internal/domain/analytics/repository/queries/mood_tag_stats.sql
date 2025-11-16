with first_query as
    (select ml.mood_log_id,
            mlmt.mood_tag_id,
            mt.name,
            date(ml.created_at) as date,
            count(mlmt.mood_tag_id) as mood_tag_id_count
     from mood_log ml
     inner join mood_log_mood_tag mlmt on ml.mood_log_id = mlmt.mood_log_id
     inner join mood_tag mt on mlmt.mood_tag_id = mt.mood_tag_id
     where ml.user_id = ?
         and date(ml.created_at) between ? and ? group  by mlmt.mood_tag_id,
                                                           mt.name,
                                                           ml.mood_log_id,
                                                           date(ml.created_at)),
     second_query as
    (select name,
            sum(mood_tag_id_count) as mood_tag_id_count ,
            (sum(mood_tag_id_count) / sum(sum(mood_tag_id_count)) over()) * 100 as percentage
     from first_query group  by mood_tag_id,
                                name)
select *
from second_query;

package database


var streaksQuery = `WITH qualifying_days AS
(
         SELECT   date,
                  created_at,
                  mood_log_id,
                  mood_rating,
                  note,
                  mood_tags,
                  mood_tag_ids,
                  daily_avg_rating,
                  daily_target_count,
                  daily_total_count,
                  daily_target_percentage
         FROM     (
                           SELECT   Date(ml.created_at) AS date,
                                    ml.created_at,
                                    ml.mood_log_id,
                                    ml.mood_rating,
                                    ml.note,
                                    group_concat(mt.NAME order BY mt.NAME separator ',')              AS mood_tags,
                                    group_concat(mt.mood_tag_id ORDER BY mt.mood_tag_id separator ',') AS mood_tag_ids,
                                    avg(ml.mood_rating) OVER (partition BY date(ml.created_at))        AS daily_avg_rating,
                                    sum(
                                    CASE
                                             WHEN mt.mood_category_id = ? THEN 1
                                             ELSE 0
                                    END) OVER (partition BY                  date(ml.created_at)) AS daily_target_count,
                                    count(mt.mood_tag_id) OVER (partition BY date(ml.created_at)) AS daily_total_count,
                                    (sum(
                                    CASE
                                             WHEN mt.mood_category_id = ? THEN 1
                                             ELSE 0
                                    END) OVER (partition BY date(ml.created_at)) * 100.0 / count(mt.mood_tag_id) OVER (partition BY date(ml.created_at))) AS daily_target_percentage
                           FROM     mood_log ml 
                           JOIN     mood_log_mood_tag mlmt
                           ON       ml.mood_log_id = mlmt.mood_log_id
                           JOIN     mood_tag mt
                           ON       mlmt.mood_tag_id = mt.mood_tag_id
                           WHERE    ml.user_id = ?
                           AND      date(ml.created_at) BETWEEN ? AND      ?
                           GROUP BY date(ml.created_at),
                                    ml.mood_log_id,
                                    ml.created_at,
                                    ml.mood_rating,
                                    ml.note) AS daily_data
         WHERE    daily_avg_rating %s ?
         AND      daily_target_percentage >= ?
         ORDER BY date,
                  created_at), second_query AS
(
         SELECT   date,
                  daily_avg_rating,
                  row_number() OVER( ORDER BY date) AS rn
         FROM     qualifying_days), third_query AS
(
         SELECT   date                              AS start_date,
                  count(*)                          AS streak_length,
                  date_add(date, interval - rn day) AS consec_groups
         FROM     second_query
         GROUP BY consec_groups), fourth_query AS
(
       SELECT *,
              date_add(start_date, interval + streak_length - 1 day) AS end_date
       FROM   third_query)
SELECT   start_date,
         end_date,
         streak_length
FROM     fourth_query
WHERE    streak_length >= 2
ORDER BY start_date;`

var daysQuery = `SELECT   date,
         created_at,
         mood_log_id,
         mood_rating,
         note,
         mood_tags,
         mood_tag_ids,
         daily_avg_rating,
         daily_target_count,
         daily_total_count,
         daily_target_percentage
FROM     (
                  SELECT   Date(ml.created_at) AS date,
                           ml.created_at,
                           ml.mood_log_id,
                           ml.mood_rating,
                           ml.note,
                           group_concat(mt.NAME order BY mt.NAME separator ', ')              AS mood_tags,
                           group_concat(mt.mood_tag_id ORDER BY mt.mood_tag_id separator ',') AS mood_tag_ids,
                           avg(ml.mood_rating) OVER (partition BY date(ml.created_at))        AS daily_avg_rating,
                           sum(
                           CASE
                                    WHEN mt.mood_category_id = ? THEN 1
                                    ELSE 0
                           END) OVER (partition BY                  date(ml.created_at)) AS daily_target_count,
                           count(mt.mood_tag_id) OVER (partition BY date(ml.created_at)) AS daily_total_count,
                           (sum(
                           CASE
                                    WHEN mt.mood_category_id = ? THEN 1
                                    ELSE 0
                           END) OVER (partition BY date(ml.created_at)) * 100.0 / count(mt.mood_tag_id) OVER (partition BY date(ml.created_at))) AS daily_target_percentage
                  FROM     mood_log ml
                  JOIN     mood_log_mood_tag mlmt
                  ON       ml.mood_log_id = mlmt.mood_log_id
                  JOIN     mood_tag mt
                  ON       mlmt.mood_tag_id = mt.mood_tag_id
                  WHERE    ml.user_id = ?
                  AND      date(ml.created_at) BETWEEN ? AND      ?
                  GROUP BY date(ml.created_at),
                           ml.mood_log_id,
                           ml.created_at,
                           ml.mood_rating,
                           ml.note ) AS daily_data
WHERE    daily_avg_rating %s ?
AND      daily_target_percentage >= ?
ORDER BY date,
         created_at;`

var stdDevQuery = `SELECT Stddev_pop(mood_rating) AS std_dev
FROM   mood_log
WHERE  user_id = ?
       AND Date(created_at) BETWEEN ? AND ?;`

var moodMovingAvgQuery = `WITH first_query
     AS (SELECT DATE(created_at) AS DATE,
                Avg(mood_rating) AS daily_avg
         FROM   mood_log
         WHERE  user_id = ?
                AND DATE(created_at) BETWEEN ? AND ?
         GROUP  BY DATE(created_at)),
     second_query
     AS (SELECT DATE,
                Avg(daily_avg)
                  OVER(
                    ORDER BY DATE ROWS BETWEEN %s preceding AND CURRENT ROW) AS
                   moving_avg
         FROM   first_query)
SELECT *
FROM   second_query
ORDER  BY DATE;`

var journalEntriesQuery = `SELECT *
FROM   mood_log
WHERE  user_id = ? and DATE(created_at) BETWEEN ? AND ?`

var moodTagFrequenciesQuery = `WITH first_query
     AS (SELECT ml.mood_log_id,
                mlmt.mood_tag_id,
                mt.NAME,
                Date(ml.created_at)    AS date,
                Count(mlmt.mood_tag_id) AS mood_tag_id_count
         FROM   mood_log ml
                INNER JOIN mood_log_mood_tag mlmt
                        ON ml.mood_log_id = mlmt.mood_log_id
                INNER JOIN mood_tag mt
                        ON mlmt.mood_tag_id = mt.mood_tag_id
         WHERE  ml.user_id = ?
                AND Date(ml.created_at) BETWEEN ? AND ?
         GROUP  BY mlmt.mood_tag_id,
                   mt.NAME,
                   ml.mood_log_id,
                   Date(ml.created_at)),
     second_query
     AS (SELECT NAME,
                Sum(mood_tag_id_count)                      AS mood_tag_id_count
                ,
                ( Sum(mood_tag_id_count) / Sum(Sum(
                  mood_tag_id_count))
                                             OVER() ) * 100 AS percentage
         FROM   first_query
         GROUP  BY mood_tag_id,
                   NAME)
SELECT *
FROM   second_query;`

var AvgMoodRatingQuery = `WITH first_query
     AS (SELECT Date(created_at) AS date,
                AVG(mood_rating) AS daily_avg_rating
         FROM   mood_log
         WHERE  user_id = ? 
                AND Date(created_at) BETWEEN ? AND ? 
         GROUP  BY Date(created_at)),
     second_query
     AS (SELECT AVG(daily_avg_rating) AS period_mood_rating_avg
         FROM   first_query)
SELECT period_mood_rating_avg
FROM   second_query;`

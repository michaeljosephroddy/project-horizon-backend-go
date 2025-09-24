package database


var streaksQueryStr = `WITH qualifying_days AS
(
         SELECT   date,
                  created_at,
                  journal_entry_id,
                  mood_rating,
                  note,
                  mood_tags,
                  mood_tag_ids,
                  daily_avg_rating,
                  daily_target_count,
                  daily_total_count,
                  daily_target_percentage
         FROM     (
                           SELECT   Date(je.created_at) AS date,
                                    je.created_at,
                                    je.journal_entry_id,
                                    je.mood_rating,
                                    je.note,
                                    group_concat(mt.NAME order BY mt.NAME separator ',')              AS mood_tags,
                                    group_concat(mt.mood_tag_id ORDER BY mt.mood_tag_id separator ',') AS mood_tag_ids,
                                    avg(je.mood_rating) OVER (partition BY date(je.created_at))        AS daily_avg_rating,
                                    sum(
                                    CASE
                                             WHEN mt.mood_category_id = ? THEN 1
                                             ELSE 0
                                    END) OVER (partition BY                  date(je.created_at)) AS daily_target_count,
                                    count(mt.mood_tag_id) OVER (partition BY date(je.created_at)) AS daily_total_count,
                                    (sum(
                                    CASE
                                             WHEN mt.mood_category_id = ? THEN 1
                                             ELSE 0
                                    END) OVER (partition BY date(je.created_at)) * 100.0 / count(mt.mood_tag_id) OVER (partition BY date(je.created_at))) AS daily_target_percentage
                           FROM     journal_entry je
                           JOIN     journal_entry_mood_tag jemt
                           ON       je.journal_entry_id = jemt.journal_entry_id
                           JOIN     mood_tag mt
                           ON       jemt.mood_tag_id = mt.mood_tag_id
                           WHERE    je.user_id = ?
                           AND      date(je.created_at) BETWEEN ? AND      ?
                           GROUP BY date(je.created_at),
                                    je.journal_entry_id,
                                    je.created_at,
                                    je.mood_rating,
                                    je.note) AS daily_data
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

var daysQueryStr = `SELECT   date,
         created_at,
         journal_entry_id,
         mood_rating,
         note,
         mood_tags,
         mood_tag_ids,
         daily_avg_rating,
         daily_target_count,
         daily_total_count,
         daily_target_percentage
FROM     (
                  SELECT   Date(je.created_at) AS date,
                           je.created_at,
                           je.journal_entry_id,
                           je.mood_rating,
                           je.note,
                           group_concat(mt.NAME order BY mt.NAME separator ', ')              AS mood_tags,
                           group_concat(mt.mood_tag_id ORDER BY mt.mood_tag_id separator ',') AS mood_tag_ids,
                           avg(je.mood_rating) OVER (partition BY date(je.created_at))        AS daily_avg_rating,
                           sum(
                           CASE
                                    WHEN mt.mood_category_id = ? THEN 1
                                    ELSE 0
                           END) OVER (partition BY                  date(je.created_at)) AS daily_target_count,
                           count(mt.mood_tag_id) OVER (partition BY date(je.created_at)) AS daily_total_count,
                           (sum(
                           CASE
                                    WHEN mt.mood_category_id = ? THEN 1
                                    ELSE 0
                           END) OVER (partition BY date(je.created_at)) * 100.0 / count(mt.mood_tag_id) OVER (partition BY date(je.created_at))) AS daily_target_percentage
                  FROM     journal_entry je
                  JOIN     journal_entry_mood_tag jemt
                  ON       je.journal_entry_id = jemt.journal_entry_id
                  JOIN     mood_tag mt
                  ON       jemt.mood_tag_id = mt.mood_tag_id
                  WHERE    je.user_id = ?
                  AND      date(je.created_at) BETWEEN ? AND      ?
                  GROUP BY date(je.created_at),
                           je.journal_entry_id,
                           je.created_at,
                           je.mood_rating,
                           je.note ) AS daily_data
WHERE    daily_avg_rating %s ?
AND      daily_target_percentage >= ?
ORDER BY date,
         created_at;`

var stdDevQueryStr = `SELECT Stddev_pop(mood_rating) AS std_dev
FROM   journal_entry
WHERE  user_id = ?
       AND Date(created_at) BETWEEN ? AND ?;`

var movingAvgQueryStr = `WITH first_query
     AS (SELECT DATE(created_at) AS DATE,
                Avg(mood_rating) AS daily_avg
         FROM   journal_entry
         WHERE  user_id = ?
                AND DATE(created_at) BETWEEN ? AND ?
         GROUP  BY DATE(created_at)),
     second_query
     AS (SELECT DATE,
                Avg(daily_avg)
                  OVER(
                    ORDER BY DATE ROWS BETWEEN 2 preceding AND CURRENT ROW) AS
                   moving_avg_3day
         FROM   first_query)
SELECT *
FROM   second_query
ORDER  BY DATE;`

var journalEntriesQueryStr = `SELECT *
FROM   journal_entry
WHERE  user_id = ? `

var moodTagFrequenciesQueryStr = `WITH first_query
     AS (SELECT je.journal_entry_id,
                jem.mood_tag_id,
                mt.NAME,
                Date(je.created_at)    AS date,
                Count(jem.mood_tag_id) AS mood_tag_id_count
         FROM   journal_entry je
                INNER JOIN journal_entry_mood_tag jem
                        ON je.journal_entry_id = jem.journal_entry_id
                INNER JOIN mood_tag mt
                        ON jem.mood_tag_id = mt.mood_tag_id
         WHERE  je.user_id = ?
                AND Date(je.created_at) BETWEEN ? AND ?
         GROUP  BY jem.mood_tag_id,
                   mt.NAME,
                   je.journal_entry_id,
                   Date(je.created_at)),
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

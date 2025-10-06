package database

var avgSleepHoursQuery = `SELECT Avg(hours_slept) AS avg_sleep_hours
FROM   sleep_log
WHERE  user_id = ?
       AND sleep_date BETWEEN ? AND ?;`

var sleepMovingAvgQuery = `WITH first_query
     AS (SELECT sleep_date AS DATE,
                Avg(hours_slept) AS avg_sleep_hours 
         FROM   sleep_log
         WHERE  user_id = ?
                AND sleep_date BETWEEN ? AND ?
	),
     second_query
     AS (SELECT DATE,
                Avg(avg_sleep_hours)
                  OVER(
                    ORDER BY DATE ROWS BETWEEN %s preceding AND CURRENT ROW) AS
                   moving_avg
         FROM   first_query)
SELECT *
FROM   second_query
ORDER  BY DATE;`

var sleepStdDevQuery = `SELECT Stddev_pop(hours_slept) AS std_dev
FROM   sleep_log
WHERE  user_id = ?
       AND sleep_date BETWEEN ? AND ?;`

var sleepQualityTagStatQuery = `WITH tag_counts AS (
    SELECT 
        sqt.name AS tag_name,
        COUNT(*) AS tag_count
    FROM sleep_log sl
    JOIN sleep_quality_tag sqt 
        ON sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
    WHERE sl.user_id = ? and sleep_date between ? and ?
    GROUP BY sqt.name
),
tag_percentages AS (
    SELECT
        tag_name,
        tag_count,
        ROUND(tag_count * 100.0 / SUM(tag_count) OVER (), 2) AS percentage
    FROM tag_counts
)
SELECT *
FROM tag_percentages
ORDER BY tag_count ASC;`

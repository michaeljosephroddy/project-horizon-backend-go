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

// var sleepQualitTagFrequenciesQuery = ``

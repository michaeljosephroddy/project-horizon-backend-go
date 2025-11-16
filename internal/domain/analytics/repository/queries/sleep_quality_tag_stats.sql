with tag_counts as
    (select sqt.name as tag_name,
            count(*) as tag_count
     from sleep_log sl
     join sleep_quality_tag sqt on sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
     where sl.user_id = ?
         and sleep_date between ? and ?
     group by sqt.name),
     tag_percentages as
    (select tag_name,
            tag_count,
            round(tag_count * 100.0 / sum(tag_count) over (), 2) as percentage
     from tag_counts)
select *
from tag_percentages
order by tag_count asc;

CREATE DATABASE IF NOT EXISTS project_horizon;
USE project_horizon;

-- Optional: Clean slate (use only in dev) - drop children first, then parents
DROP TABLE IF EXISTS mood_log_mood_tag;
DROP TABLE IF EXISTS medication_log_item;  -- Added this line!
DROP TABLE IF EXISTS medication_log;
DROP TABLE IF EXISTS user_medication;
DROP TABLE IF EXISTS mood_log;
DROP TABLE IF EXISTS sleep_log;
DROP TABLE IF EXISTS mood_tag;
DROP TABLE IF EXISTS mood_category;
DROP TABLE IF EXISTS medication;
DROP TABLE IF EXISTS sleep_quality_tag;
DROP TABLE IF EXISTS user;

-- User table
CREATE TABLE IF NOT EXISTS user (
    user_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
);

-- Medications
CREATE TABLE IF NOT EXISTS medication (
    medication_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mood categories
CREATE TABLE IF NOT EXISTS mood_category (
    mood_category_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    description VARCHAR(100)
);

-- Mood tags
CREATE TABLE IF NOT EXISTS mood_tag (
    mood_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    mood_category_id BIGINT UNSIGNED NOT NULL,
    CONSTRAINT fk_mood_tag_category FOREIGN KEY (mood_category_id) REFERENCES mood_category(mood_category_id),
    INDEX idx_category (mood_category_id)
);

-- User medications for med history
CREATE TABLE IF NOT EXISTS user_medication (
    user_medication_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    medication_id BIGINT UNSIGNED NOT NULL,
    dosage VARCHAR(50),
    start_date DATE NOT NULL,
    end_date DATE,
    stopped TINYINT(1) NOT NULL DEFAULT 0,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_medication_user FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_medication_med FOREIGN KEY (medication_id) REFERENCES medication(medication_id) ON DELETE CASCADE,
    INDEX idx_user_med (user_id, medication_id),
    INDEX idx_start_date (start_date)
);

-- Daily adherence log (header/parent record)
CREATE TABLE IF NOT EXISTS medication_log (
    medication_log_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    taken_at TIMESTAMP NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_medication_log_user FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    INDEX idx_taken_at (taken_at),
    INDEX idx_user_taken (user_id, taken_at)
);

-- Join table for medications taken in each log entry
CREATE TABLE IF NOT EXISTS medication_log_item (
    medication_log_item_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    medication_log_id BIGINT UNSIGNED NOT NULL,
    medication_id BIGINT UNSIGNED NOT NULL,
    dosage VARCHAR(50) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_log_item_log FOREIGN KEY (medication_log_id) REFERENCES medication_log(medication_log_id) ON DELETE CASCADE,
    CONSTRAINT fk_log_item_med FOREIGN KEY (medication_id) REFERENCES medication(medication_id) ON DELETE CASCADE,
    INDEX idx_log_id (medication_log_id),
    INDEX idx_medication_id (medication_id)
);

-- Mood log entries
CREATE TABLE IF NOT EXISTS mood_log (
    mood_log_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    mood_rating INT NOT NULL CHECK (mood_rating BETWEEN 1 AND 10),
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_mood_log_user FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    INDEX idx_user_created (user_id, created_at),
    INDEX idx_created_at (created_at)
);

-- Mood log mood tags join table
CREATE TABLE IF NOT EXISTS mood_log_mood_tag (
    mood_log_mood_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    mood_log_id BIGINT UNSIGNED NOT NULL,
    mood_tag_id BIGINT UNSIGNED NOT NULL,
    CONSTRAINT fk_mlmt_log FOREIGN KEY (mood_log_id) REFERENCES mood_log(mood_log_id) ON DELETE CASCADE,
    CONSTRAINT fk_mlmt_tag FOREIGN KEY (mood_tag_id) REFERENCES mood_tag(mood_tag_id) ON DELETE CASCADE,
    INDEX idx_mood_log (mood_log_id),
    INDEX idx_mood_tag (mood_tag_id),
    UNIQUE KEY unique_log_tag (mood_log_id, mood_tag_id)
);

-- The different sleep quality tags
CREATE TABLE IF NOT EXISTS sleep_quality_tag (
    sleep_quality_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Track Sleep Entries
CREATE TABLE IF NOT EXISTS sleep_log (
    sleep_log_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    hours_slept DECIMAL(4,2) NOT NULL CHECK (hours_slept >= 0 AND hours_slept <= 24),
    sleep_quality_tag_id BIGINT UNSIGNED NOT NULL,
    note TEXT,
    sleep_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_sleep_log_user FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_sleep_log_quality FOREIGN KEY (sleep_quality_tag_id) REFERENCES sleep_quality_tag(sleep_quality_tag_id) ON DELETE CASCADE,
    INDEX idx_user_date (user_id, sleep_date),
    INDEX idx_sleep_date (sleep_date)
);

-- Reset auto-increments
ALTER TABLE user AUTO_INCREMENT = 1;
ALTER TABLE mood_category AUTO_INCREMENT = 1;
ALTER TABLE mood_tag AUTO_INCREMENT = 1;
ALTER TABLE mood_log AUTO_INCREMENT = 1;
ALTER TABLE medication AUTO_INCREMENT = 1;
ALTER TABLE user_medication AUTO_INCREMENT = 1;
ALTER TABLE medication_log AUTO_INCREMENT = 1;
ALTER TABLE medication_log_item AUTO_INCREMENT = 1;
ALTER TABLE mood_log_mood_tag AUTO_INCREMENT = 1;
ALTER TABLE sleep_log AUTO_INCREMENT = 1;
ALTER TABLE sleep_quality_tag AUTO_INCREMENT = 1;

-- DB user
CREATE USER IF NOT EXISTS 'demouser'@'localhost' IDENTIFIED BY 'demopassword';
GRANT ALL PRIVILEGES ON project_horizon.* TO 'demouser'@'localhost';
FLUSH PRIVILEGES;

-- ========================================
-- SEED DATA
-- ========================================

-- Mood categories
INSERT INTO mood_category (name, description) VALUES
('positive', 'Positive emotions and feelings'),
('negative', 'Challenging or difficult emotions'),
('neutral', 'Neutral or mixed emotional states'),
('energy', 'Energy and physical state related moods'),
('clinical', 'Clinical mood states related to bipolar disorder');

-- Mood tags
INSERT INTO mood_tag (name, mood_category_id) VALUES
-- Positive
('Happy', 1), ('Excited', 1), ('Calm', 1), ('Grateful', 1), ('Confident', 1),
-- Negative
('Sad', 2), ('Anxious', 2), ('Angry', 2), ('Frustrated', 2), ('Lonely', 2),
-- Neutral
('Content', 3), ('Restless', 3), ('Confused', 3), ('Bored', 3),
-- Energy
('Energetic', 4), ('Tired', 4),
-- Clinical
('Manic', 5), ('Hypomanic', 5), ('Depressed', 5), ('Mixed State', 5), ('Irritable', 5);

-- Insert users
INSERT INTO user (email, password_hash, created_at) VALUES
('alice@example.com', '$2y$10$example.hash.1', '2025-07-15 10:00:00'),
('bob@example.com', '$2y$10$example.hash.2', '2025-07-20 14:30:00'),
('carol@example.com', '$2y$10$example.hash.3', '2025-07-25 09:15:00');

-- Medications
INSERT INTO medication (name, description) VALUES
('Sertraline', 'SSRI antidepressant used to treat depression, anxiety, OCD, and PTSD'),
('Fluoxetine', 'SSRI antidepressant commonly known as Prozac'),
('Escitalopram', 'SSRI antidepressant used for depression and generalized anxiety disorder'),
('Bupropion', 'Atypical antidepressant that can help with depression and smoking cessation'),
('Venlafaxine', 'SNRI antidepressant for depression and anxiety disorders'),
('Lithium', 'Mood stabilizer primarily used to treat bipolar disorder'),
('Lamotrigine', 'Mood stabilizer used for bipolar disorder and seizure prevention'),
('Quetiapine', 'Atypical antipsychotic used for bipolar disorder, schizophrenia, and depression'),
('Aripiprazole', 'Atypical antipsychotic for schizophrenia, bipolar disorder, and depression'),
('Valproate', 'Mood stabilizer used for bipolar disorder, seizures, and migraine prevention');

-- Link users/medications
INSERT INTO user_medication (user_id, medication_id, dosage, start_date, stopped) VALUES
(1, 1, '50mg', '2025-01-15', 0),
(1, 6, '300mg', '2025-02-01', 0),
(2, 2, '20mg', '2025-03-10', 0),
(2, 7, '100mg', '2025-03-15', 0),
(3, 3, '10mg', '2025-04-01', 0),
(3, 8, '200mg', '2025-04-10', 0);

-- Sleep quality tags
INSERT INTO sleep_quality_tag (name, description) VALUES
('Excellent', 'Slept deeply, woke up refreshed and energized'),
('Good', 'Slept well with minimal disruptions'),
('Fair', 'Adequate sleep but not fully restorative'),
('Poor', 'Restless sleep with frequent waking'),
('Very Poor', 'Minimal sleep, exhausted upon waking');

-- Sample sleep logs for August 2025
INSERT INTO sleep_log (user_id, hours_slept, sleep_quality_tag_id, sleep_date, note) VALUES
-- Alice's sleep logs
(1, 7.5, 2, '2025-08-01', 'Slept well, ready for the day'),
(1, 8.0, 1, '2025-08-02', 'Amazing sleep after great day'),
(1, 7.0, 2, '2025-08-03', 'Good rest'),
(1, 8.5, 1, '2025-08-04', 'Perfect sleep after hiking'),
(1, 6.5, 3, '2025-08-05', 'Monday anxiety affected sleep'),
-- Bob's sleep logs
(2, 6.0, 3, '2025-08-01', 'Decent sleep'),
(2, 5.5, 4, '2025-08-02', 'Restless night'),
(2, 7.0, 2, '2025-08-03', 'Better sleep after park walk'),
(2, 7.5, 2, '2025-08-04', 'Good rest'),
(2, 7.0, 2, '2025-08-05', 'Solid sleep'),
-- Carol's sleep logs
(3, 6.5, 3, '2025-08-01', 'Okay sleep, some anxiety'),
(3, 7.0, 2, '2025-08-02', 'Better sleep after good conversation'),
(3, 6.0, 3, '2025-08-03', 'Average sleep'),
(3, 8.0, 1, '2025-08-04', 'Excellent sleep after family day'),
(3, 5.5, 4, '2025-08-05', 'Stressed about deadlines');

-- Journal entries for August 2025 (31 days)
-- Cycling through user_id 1, 2, and 3
INSERT INTO mood_log (user_id, mood_rating, note, created_at) VALUES
-- August 1, 2025
(1, 7, 'Started the month feeling optimistic. Work project is going well and I had a great workout this morning.', '2025-08-01 08:30:00'),
(2, 5, 'Feeling neutral today. Nothing particularly good or bad happened. Just a regular Thursday.', '2025-08-01 19:45:00'),
(3, 6, 'Had some anxiety this morning but it passed after my meditation session. Grateful for small victories.', '2025-08-01 21:15:00'),

-- August 2, 2025
(1, 8, 'Excellent day! Completed a major milestone at work and celebrated with friends. Feeling accomplished.', '2025-08-02 22:00:00'),
(2, 4, 'Woke up feeling a bit down. Weather is gloomy and it is affecting my mood. Need to find indoor activities.', '2025-08-02 10:20:00'),
(3, 7, 'Really good day overall. Connected with an old friend and we had a wonderful conversation over coffee.', '2025-08-02 16:30:00'),

-- August 3, 2025
(1, 6, 'Weekend vibes are kicking in. Feeling relaxed but also productive. Organized my living space.', '2025-08-03 14:00:00'),
(2, 6, 'Better than yesterday. Went for a walk in the park and it lifted my spirits. Nature therapy works.', '2025-08-03 17:45:00'),
(3, 5, 'Average day. Work was busy but manageable. Looking forward to the weekend to recharge.', '2025-08-03 20:30:00'),

-- August 4, 2025
(1, 9, 'Amazing Sunday! Went hiking with friends and saw the most beautiful sunrise. Feeling so grateful.', '2025-08-04 19:00:00'),
(2, 7, 'Had a productive day. Caught up on reading and tried a new recipe. Small pleasures make a difference.', '2025-08-04 21:30:00'),
(3, 8, 'Wonderful family day. Had a barbecue with relatives and laughed until my sides hurt. Pure joy.', '2025-08-04 22:45:00'),

-- August 5, 2025
(1, 5, 'Monday blues hit hard. Back to work after a great weekend. Need to find better work-life balance.', '2025-08-05 18:00:00'),
(2, 6, 'Decent start to the week. Had a good meeting at work and received positive feedback on my project.', '2025-08-05 16:15:00'),
(3, 4, 'Feeling overwhelmed with deadlines this week. Need to prioritize and take things one step at a time.', '2025-08-05 20:00:00'),

-- August 6, 2025
(1, 7, 'Much better today. Found my rhythm at work and made good progress. Evening yoga helped center me.', '2025-08-06 21:00:00'),
(2, 5, 'Neutral day. Nothing exciting but nothing terrible either. Sometimes average is perfectly fine.', '2025-08-06 19:30:00'),
(3, 6, 'Managed my stress better today. Used breathing exercises between tasks and it really helped.', '2025-08-06 18:45:00'),

-- August 7, 2025
(1, 8, 'Great energy today! Completed all my tasks and even helped a colleague with their project. Teamwork!', '2025-08-07 17:30:00'),
(2, 7, 'Feeling more positive. Had lunch with a friend and we shared some good laughs. Social connection matters.', '2025-08-07 20:15:00'),
(3, 7, 'Good day at work and evening walk with my partner. Simple pleasures are the best.', '2025-08-07 21:45:00'),

-- August 8, 2025
(1, 6, 'Thursday feeling. Ready for the weekend but still focused on finishing strong this week.', '2025-08-08 18:30:00'),
(2, 8, 'Surprisingly good day! Received unexpected praise from my manager and treated myself to a nice dinner.', '2025-08-08 22:00:00'),
(3, 5, 'Tired today but pushed through. Looking forward to resting this weekend and recharging.', '2025-08-08 19:00:00'),

-- August 9, 2025
(1, 9, 'FRIDAY! And what a fantastic one. Finished a big project and the weekend is here. Time to celebrate!', '2025-08-09 17:00:00'),
(2, 6, 'End of week energy. Not the best day but glad it is Friday. Weekend plans are helping my mood.', '2025-08-09 18:45:00'),
(3, 7, 'Good end to the work week. Accomplished my goals and feeling prepared for next week. Balance is key.', '2025-08-09 20:30:00'),

-- August 10, 2025
(1, 7, 'Relaxing Saturday. Slept in, read a good book, and cooked a nice meal. Self-care Saturday success.', '2025-08-10 19:45:00'),
(2, 8, 'Wonderful Saturday! Went to a farmers market and tried new foods. Exploring new things energizes me.', '2025-08-10 21:00:00'),
(3, 6, 'Quiet weekend day. Did some household chores and watched movies. Sometimes low-key is perfect.', '2025-08-10 22:15:00'),

-- August 11, 2025
(1, 8, 'Perfect Sunday! Brunch with friends, afternoon in the park, and evening movie. Ideal weekend conclusion.', '2025-08-11 21:30:00'),
(2, 7, 'Nice Sunday. Called family and caught up with everyone. Family connections always boost my mood.', '2025-08-11 20:00:00'),
(3, 7, 'Good weekend wrap-up. Prepared for the week ahead and spent quality time with loved ones.', '2025-08-11 19:30:00'),

-- August 12, 2025
(1, 5, 'Monday again. Feeling the weekend withdrawal but trying to start the week positively. Coffee helps.', '2025-08-12 09:00:00'),
(2, 6, 'Decent Monday start. New week, new opportunities. Trying to maintain a positive mindset.', '2025-08-12 17:45:00'),
(3, 4, 'Monday stress is real. Lots on my plate this week but breaking it down into manageable chunks.', '2025-08-12 18:30:00'),

-- August 13, 2025
(1, 7, 'Tuesday improvement! Getting into the week is rhythm. Productive day and good team collaboration.', '2025-08-13 18:00:00'),
(2, 6, 'Steady Tuesday. Making progress on projects and feeling more confident about the week ahead.', '2025-08-13 19:15:00'),
(3, 6, 'Better day today. Tackled some challenging tasks and feeling more in control of my workload.', '2025-08-13 20:45:00'),

-- August 14, 2025
(1, 8, 'Hump day excellence! Mid-week and feeling strong. Great meeting outcomes and positive energy.', '2025-08-14 17:15:00'),
(2, 7, 'Wednesday win! Solved a problem that had been bothering me for days. Breakthrough moments feel amazing.', '2025-08-14 21:00:00'),
(3, 7, 'Solid mid-week day. Feeling more balanced and confident in handling work challenges.', '2025-08-14 19:00:00'),

-- August 15, 2025
(1, 6, 'Thursday thoughts. Almost to the weekend but staying present. Completed important tasks today.', '2025-08-15 18:45:00'),
(2, 5, 'Neutral Thursday. Not great, not bad. Sometimes these steady days are exactly what I need.', '2025-08-15 20:30:00'),
(3, 8, 'Excellent Thursday! Received great feedback and feeling accomplished. Hard work is paying off.', '2025-08-15 19:45:00'),

-- August 16, 2025
(1, 9, 'Friday fantastic! Week accomplished, weekend ahead, and celebrating small and big victories.', '2025-08-16 17:30:00'),
(2, 7, 'Good Friday finish. Ready for weekend adventures and feeling grateful for a productive week.', '2025-08-16 18:15:00'),
(3, 6, 'End of week exhaustion but also satisfaction. Accomplished goals and ready to rest and recharge.', '2025-08-16 21:00:00'),

-- August 17, 2025
(1, 8, 'Saturday bliss! Morning jog, afternoon with friends, evening relaxation. Perfect balance achieved.', '2025-08-17 20:00:00'),
(2, 8, 'Wonderful Saturday! Tried a new hiking trail and discovered a beautiful viewpoint. Adventure therapy.', '2025-08-17 19:30:00'),
(3, 7, 'Nice Saturday. Balanced productivity with relaxation. Got things done but also took care of myself.', '2025-08-17 21:15:00'),

-- August 18, 2025
(1, 7, 'Sunday funday! Good mix of activities and rest. Feeling recharged for the upcoming week.', '2025-08-18 19:45:00'),
(2, 6, 'Quiet Sunday. Sometimes the best weekends are the ones where you do not plan much. Peace.', '2025-08-18 20:45:00'),
(3, 8, 'Excellent Sunday! Quality time with family and friends. Feeling loved and supported.', '2025-08-18 21:30:00'),

-- August 19, 2025
(1, 6, 'Monday motivation is building. New week, fresh start, and feeling prepared for whatever comes.', '2025-08-19 08:45:00'),
(2, 5, 'Another Monday. Trying to approach it with curiosity instead of dread. Small mindset shifts help.', '2025-08-19 17:30:00'),
(3, 5, 'Monday manageable. Not excited but not dreading it either. Finding the middle ground.', '2025-08-19 18:00:00'),

-- August 20, 2025
(1, 7, 'Tuesday triumph! Good progress on projects and positive interactions with colleagues. Momentum building.', '2025-08-20 18:30:00'),
(2, 7, 'Better Tuesday. Had an inspiring conversation that gave me new ideas and energy. Connection matters.', '2025-08-20 20:00:00'),
(3, 6, 'Steady Tuesday progress. Taking things step by step and feeling more confident each day.', '2025-08-20 19:15:00'),

-- August 21, 2025
(1, 8, 'Wednesday wonderful! Mid-week and feeling excellent. Great flow state at work and evening workout.', '2025-08-21 21:00:00'),
(2, 6, 'Decent Wednesday. Nothing extraordinary but solid progress. Sometimes consistency is the goal.', '2025-08-21 19:45:00'),
(3, 7, 'Good Wednesday energy. Tackled challenging tasks and feeling accomplished. Building confidence.', '2025-08-21 18:45:00'),

-- August 22, 2025
(1, 5, 'Thursday thoughts turning to weekend. Energy dipping but pushing through. Almost there!', '2025-08-22 17:45:00'),
(2, 8, 'Surprising Thursday high! Unexpected good news brightened my whole day. Grateful for pleasant surprises.', '2025-08-22 20:30:00'),
(3, 6, 'Thursday steady. Maintaining good momentum and looking forward to Friday accomplishments.', '2025-08-22 19:30:00'),

-- August 23, 2025
(1, 9, 'Friday celebration! Incredible week completion and weekend adventures awaiting. Life is good!', '2025-08-23 17:00:00'),
(2, 7, 'Happy Friday! Week had ups and downs but ending strong. Ready for weekend restoration.', '2025-08-23 18:30:00'),
(3, 8, 'Fantastic Friday finish! Achieved weekly goals and feeling proud of the progress made.', '2025-08-23 19:00:00'),

-- August 24, 2025
(1, 7, 'Saturday satisfaction! Good balance of productivity and relaxation. Exactly what weekends should be.', '2025-08-24 20:15:00'),
(2, 8, 'Amazing Saturday adventure! Explored a new part of the city and discovered hidden gems.', '2025-08-24 21:45:00'),
(3, 6, 'Calm Saturday. Low-key activities and quality time at home. Sometimes simple is best.', '2025-08-24 19:45:00'),

-- August 25, 2025
(1, 8, 'Sunday success! Perfect end to the weekend with good food, great company, and relaxation.', '2025-08-25 20:30:00'),
(2, 7, 'Nice Sunday wrap-up. Prepared for the week while still enjoying weekend vibes. Balance achieved.', '2025-08-25 19:15:00'),
(3, 7, 'Good Sunday conclusion. Ready for a new week with optimism and energy restored.', '2025-08-25 21:00:00'),

-- August 26, 2025
(1, 6, 'Monday mindset improving. Starting to see Mondays as opportunities rather than obstacles.', '2025-08-26 09:15:00'),
(2, 5, 'Standard Monday. Not thrilled but not terrible. Focusing on small wins throughout the day.', '2025-08-26 18:00:00'),
(3, 6, 'Monday momentum building. Good energy and clear priorities for the week ahead.', '2025-08-26 17:45:00'),

-- August 27, 2025
(1, 7, 'Tuesday productivity peak! Everything clicked today and made excellent progress on key projects.', '2025-08-27 18:15:00'),
(2, 6, 'Decent Tuesday development. Steady progress and maintaining positive attitude despite challenges.', '2025-08-27 19:30:00'),
(3, 7, 'Strong Tuesday performance. Feeling confident and capable. Good rhythm established for the week.', '2025-08-27 20:00:00'),

-- August 28, 2025
(1, 8, 'Wednesday winner! Mid-week excellence with great achievements and positive team interactions.', '2025-08-28 17:30:00'),
(2, 7, 'Wednesday breakthrough! Solved a complex problem and feeling intellectually satisfied.', '2025-08-28 21:15:00'),
(3, 6, 'Steady Wednesday progress. Consistent effort and maintaining good work-life balance.', '2025-08-28 19:00:00'),

-- August 29, 2025
(1, 6, 'Thursday transition. Moving toward weekend mode but staying focused on finishing strong.', '2025-08-29 18:45:00'),
(2, 8, 'Excellent Thursday! Unexpected positive developments made this day special. Grateful for surprises.', '2025-08-29 20:45:00'),
(3, 7, 'Good Thursday execution. Accomplished daily goals and feeling prepared for Friday finals.', '2025-08-29 19:30:00'),

-- August 30, 2025
(1, 9, 'Friday finale fantastic! Amazing end to August with accomplishments and weekend excitement ahead!', '2025-08-30 17:15:00'),
(2, 7, 'Happy Friday and month conclusion! August had its challenges but overall positive growth.', '2025-08-30 18:30:00'),
(3, 8, 'Wonderful Friday wrap-up! Reflecting on August achievements and looking forward to September.', '2025-08-30 20:15:00'),

-- August 31, 2025
(1, 8, 'Saturday celebration! Last day of August spent perfectly with reflection and anticipation.', '2025-08-31 19:30:00'),
(2, 7, 'August conclusion satisfaction. Good month overall with growth, challenges overcome, and joy found.', '2025-08-31 21:00:00'),
(3, 7, 'Perfect August ending! Grateful for the experiences and ready for new adventures in September.', '2025-08-31 20:45:00');

-- Mood tag associations for journal entries
INSERT INTO mood_log_mood_tag (mood_log_id, mood_tag_id) VALUES
-- August 1, 2025
(1, 1), (1, 15), -- Alice: Happy, Energetic (optimistic, great workout)
(2, 11), (2, 14), -- Bob: Content, Bored (neutral, regular day)
(3, 7), (3, 4), -- Carol: Anxious, Grateful (anxiety but grateful)

-- August 2, 2025
(4, 1), (4, 5), (4, 2), -- Alice: Happy, Confident, Excited (excellent day, accomplished)
(5, 6), (5, 16), -- Bob: Sad, Tired (down, gloomy weather)
(6, 1), (6, 11), -- Carol: Happy, Content (wonderful conversation)

-- August 3, 2025
(7, 3), (7, 11), -- Alice: Calm, Content (relaxed, productive)
(8, 1), (8, 3), -- Bob: Happy, Calm (spirits lifted by walk)
(9, 11), (9, 16), -- Carol: Content, Tired (busy but manageable)

-- August 4, 2025
(10, 1), (10, 4), (10, 15), -- Alice: Happy, Grateful, Energetic (amazing sunrise hike)
(11, 11), (11, 1), -- Bob: Content, Happy (productive day, small pleasures)
(12, 1), (12, 2), -- Carol: Happy, Excited (family barbecue, pure joy)

-- August 5, 2025
(13, 6), (13, 16), -- Alice: Sad, Tired (Monday blues)
(14, 11), (14, 5), -- Bob: Content, Confident (positive feedback)
(15, 12), (15, 7), -- Carol: Restless, Anxious (overwhelmed with deadlines)

-- August 6, 2025
(16, 1), (16, 3), -- Alice: Happy, Calm (found rhythm, yoga)
(17, 11), -- Bob: Content (neutral day)
(18, 3), (18, 11), -- Carol: Calm, Content (managed stress better)

-- August 7, 2025
(19, 1), (19, 15), (19, 2), -- Alice: Happy, Energetic, Excited (great energy, teamwork)
(20, 1), (20, 11), -- Bob: Happy, Content (positive, good laughs)
(21, 1), (21, 3), -- Carol: Happy, Calm (good day, simple pleasures)

-- August 8, 2025
(22, 11), (22, 16), -- Alice: Content, Tired (ready for weekend)
(23, 1), (23, 2), -- Bob: Happy, Excited (unexpected praise)
(24, 16), (24, 11), -- Carol: Tired, Content (tired but pushed through)

-- August 9, 2025
(25, 1), (25, 2), (25, 15), -- Alice: Happy, Excited, Energetic (FRIDAY celebration)
(26, 11), (26, 1), -- Bob: Content, Happy (glad it's Friday)
(27, 1), (27, 5), -- Carol: Happy, Confident (accomplished goals)

-- August 10, 2025
(28, 3), (28, 11), -- Alice: Calm, Content (relaxing Saturday)
(29, 1), (29, 2), (29, 15), -- Bob: Happy, Excited, Energetic (farmers market adventure)
(30, 3), (30, 11), -- Carol: Calm, Content (quiet, low-key)

-- August 11, 2025
(31, 1), (31, 11), -- Alice: Happy, Content (perfect Sunday)
(32, 1), (32, 3), -- Bob: Happy, Calm (nice family calls)
(33, 1), (33, 11), -- Carol: Happy, Content (quality time)

-- August 12, 2025
(34, 6), (34, 16), -- Alice: Sad, Tired (Monday withdrawal)
(35, 11), (35, 5), -- Bob: Content, Confident (positive mindset)
(36, 12), (36, 7), -- Carol: Restless, Anxious (Monday stress)

-- August 13, 2025
(37, 1), (37, 5), -- Alice: Happy, Confident (productive, good collaboration)
(38, 11), (38, 5), -- Bob: Content, Confident (making progress)
(39, 5), (39, 11), -- Carol: Confident, Content (more in control)

-- August 14, 2025
(40, 1), (40, 15), (40, 2), -- Alice: Happy, Energetic, Excited (hump day excellence)
(41, 1), (41, 2), -- Bob: Happy, Excited (breakthrough moment)
(42, 11), (42, 5), -- Carol: Content, Confident (balanced)

-- August 15, 2025
(43, 11), (43, 5), -- Alice: Content, Confident (completed tasks)
(44, 11), -- Bob: Content (neutral day)
(45, 1), (45, 2), -- Carol: Happy, Excited (great feedback)

-- August 16, 2025
(46, 1), (46, 2), (46, 15), -- Alice: Happy, Excited, Energetic (Friday fantastic)
(47, 1), (47, 4), -- Bob: Happy, Grateful (productive week)
(48, 16), (48, 11), -- Carol: Tired, Content (exhausted but satisfied)

-- August 17, 2025
(49, 1), (49, 3), -- Alice: Happy, Calm (Saturday bliss, perfect balance)
(50, 1), (50, 2), (50, 15), -- Bob: Happy, Excited, Energetic (new hiking trail)
(51, 11), (51, 3), -- Carol: Content, Calm (balanced productivity)

-- August 18, 2025
(52, 1), (52, 3), -- Alice: Happy, Calm (Sunday funday, recharged)
(53, 3), (53, 11), -- Bob: Calm, Content (quiet Sunday, peace)
(54, 1), (54, 4), -- Carol: Happy, Grateful (quality time, supported)

-- August 19, 2025
(55, 5), (55, 11), -- Alice: Confident, Content (motivated, prepared)
(56, 11), (56, 13), -- Bob: Content, Confused (trying new mindset)
(57, 11), -- Carol: Content (middle ground)

-- August 20, 2025
(58, 1), (58, 5), -- Alice: Happy, Confident (triumph, momentum)
(59, 1), (59, 2), -- Bob: Happy, Excited (inspiring conversation)
(60, 5), (60, 11), -- Carol: Confident, Content (step by step progress)

-- August 21, 2025
(61, 1), (61, 15), -- Alice: Happy, Energetic (wonderful, flow state)
(62, 11), (62, 5), -- Bob: Content, Confident (solid progress)
(63, 1), (63, 5), -- Carol: Happy, Confident (accomplished)

-- August 22, 2025
(64, 16), (64, 12), -- Alice: Tired, Restless (energy dipping)
(65, 1), (65, 2), -- Bob: Happy, Excited (unexpected good news)
(66, 11), (66, 5), -- Carol: Content, Confident (steady momentum)

-- August 23, 2025
(67, 1), (67, 2), (67, 15), -- Alice: Happy, Excited, Energetic (Friday celebration)
(68, 1), (68, 11), -- Bob: Happy, Content (ending strong)
(69, 1), (69, 5), -- Carol: Happy, Confident (achieved goals)

-- August 24, 2025
(70, 11), (70, 3), -- Alice: Content, Calm (satisfaction, balance)
(71, 1), (71, 2), (71, 15), -- Bob: Happy, Excited, Energetic (amazing adventure)
(72, 3), (72, 11), -- Carol: Calm, Content (low-key, simple)

-- August 25, 2025
(73, 1), (73, 11), -- Alice: Happy, Content (perfect end, good company)
(74, 11), (74, 3), -- Bob: Content, Calm (prepared, balance)
(75, 1), (75, 2), -- Carol: Happy, Excited (optimism, energy restored)

-- August 26, 2025
(76, 5), (76, 11), -- Alice: Confident, Content (mindset improving)
(77, 11), -- Bob: Content (standard Monday)
(78, 15), (78, 5), -- Carol: Energetic, Confident (momentum building)

-- August 27, 2025
(79, 1), (79, 5), -- Alice: Happy, Confident (productivity peak)
(80, 11), (80, 5), -- Bob: Content, Confident (steady progress)
(81, 1), (81, 5), -- Carol: Happy, Confident (strong performance)

-- August 28, 2025
(82, 1), (82, 15), (82, 2), -- Alice: Happy, Energetic, Excited (Wednesday winner)
(83, 1), (83, 5), -- Bob: Happy, Confident (breakthrough)
(84, 11), (84, 5), -- Carol: Content, Confident (steady progress)

-- August 29, 2025
(85, 11), (85, 5), -- Alice: Content, Confident (focused finish)
(86, 1), (86, 2), -- Bob: Happy, Excited (positive developments)
(87, 1), (87, 5), -- Carol: Happy, Confident (accomplished goals)

-- August 30, 2025
(88, 1), (88, 2), (88, 15), -- Alice: Happy, Excited, Energetic (Friday finale fantastic)
(89, 1), (89, 11), -- Bob: Happy, Content (positive growth)
(90, 1), (90, 4), -- Carol: Happy, Grateful (reflecting on achievements)

-- August 31, 2025
(91, 1), (91, 4), -- Alice: Happy, Grateful (celebration, reflection)
(92, 11), (92, 4), -- Bob: Content, Grateful (good month overall)
(93, 1), (93, 4); -- Carol: Happy, Grateful (perfect ending)

-- Mock data for medication_log (August 2025)
-- User 1: Takes Sertraline (50mg) and Lithium (300mg)
-- User 2: Takes Fluoxetine (20mg) and Lamotrigine (100mg)
-- User 3: Takes Escitalopram (10mg) and Quetiapine (200mg)

-- User 1 Logs
INSERT INTO medication_log (user_id, taken_at, note) VALUES
(1, '2025-08-01 08:30:00', 'Morning dose, feeling good'),
(1, '2025-08-02 09:15:00', NULL),
(1, '2025-08-03 08:45:00', 'Slept well last night'),
(1, '2025-08-04 14:20:00', 'Forgot morning dose, taking now'),
(1, '2025-08-05 08:00:00', NULL),
(1, '2025-08-06 09:30:00', 'Feeling more stable'),
(1, '2025-08-07 08:15:00', NULL),
(1, '2025-08-08 19:00:00', 'Evening dose, busy day'),
(1, '2025-08-09 08:45:00', NULL),
(1, '2025-08-10 09:00:00', 'Weekend, relaxed morning'),
(1, '2025-08-11 10:30:00', NULL),
(1, '2025-08-12 08:30:00', 'Back to routine'),
(1, '2025-08-13 08:15:00', NULL),
(1, '2025-08-14 14:45:00', 'Afternoon dose'),
(1, '2025-08-15 08:00:00', 'Feeling energized'),
(1, '2025-08-16 09:15:00', NULL),
(1, '2025-08-17 08:30:00', 'Good mood today'),
(1, '2025-08-18 15:30:00', 'Late dose, hectic schedule'),
(1, '2025-08-19 08:45:00', NULL),
(1, '2025-08-20 09:00:00', NULL),
(1, '2025-08-21 08:15:00', 'Mild headache this morning'),
(1, '2025-08-22 08:30:00', NULL),
(1, '2025-08-23 19:30:00', 'Evening dose after work'),
(1, '2025-08-24 09:30:00', 'Weekend morning'),
(1, '2025-08-25 10:00:00', NULL),
(1, '2025-08-26 08:00:00', 'Starting the week strong'),
(1, '2025-08-27 08:45:00', NULL),
(1, '2025-08-28 09:15:00', 'Feeling balanced'),
(1, '2025-08-29 14:00:00', 'Afternoon dose'),
(1, '2025-08-30 08:30:00', NULL),
(1, '2025-08-31 09:00:00', 'End of month check-in');

-- User 2 Logs
INSERT INTO medication_log (user_id, taken_at, note) VALUES
(2, '2025-08-01 07:45:00', 'Morning routine'),
(2, '2025-08-02 08:30:00', NULL),
(2, '2025-08-03 13:15:00', 'Took with lunch'),
(2, '2025-08-04 07:30:00', 'Early start today'),
(2, '2025-08-05 08:00:00', NULL),
(2, '2025-08-06 08:15:00', 'Feeling calm'),
(2, '2025-08-07 18:45:00', 'Evening dose, forgot earlier'),
(2, '2025-08-08 07:45:00', NULL),
(2, '2025-08-09 08:30:00', 'Good sleep quality'),
(2, '2025-08-10 09:45:00', 'Lazy Saturday morning'),
(2, '2025-08-11 10:15:00', NULL),
(2, '2025-08-12 07:30:00', 'Back to work schedule'),
(2, '2025-08-13 08:00:00', NULL),
(2, '2025-08-14 08:15:00', 'Mood is stable'),
(2, '2025-08-15 14:30:00', 'Midday dose'),
(2, '2025-08-16 07:45:00', NULL),
(2, '2025-08-17 08:30:00', 'Feeling productive'),
(2, '2025-08-18 08:00:00', NULL),
(2, '2025-08-19 19:15:00', 'Late evening dose'),
(2, '2025-08-20 08:15:00', NULL),
(2, '2025-08-21 07:30:00', 'Early morning'),
(2, '2025-08-22 08:45:00', 'Slight nausea but manageable'),
(2, '2025-08-23 08:00:00', NULL),
(2, '2025-08-24 10:30:00', 'Weekend brunch time'),
(2, '2025-08-25 09:00:00', NULL),
(2, '2025-08-26 07:45:00', 'Monday morning'),
(2, '2025-08-27 08:30:00', NULL),
(2, '2025-08-28 08:00:00', 'Anxiety levels low'),
(2, '2025-08-29 15:00:00', 'Afternoon dose, busy day'),
(2, '2025-08-30 08:15:00', NULL),
(2, '2025-08-31 08:30:00', 'Month went well');

-- User 3 Logs
INSERT INTO medication_log (user_id, taken_at, note) VALUES
(3, '2025-08-01 09:00:00', NULL),
(3, '2025-08-01 21:00:00', 'Evening Quetiapine dose'),
(3, '2025-08-02 08:45:00', 'Morning meds'),
(3, '2025-08-03 09:15:00', NULL),
(3, '2025-08-03 20:30:00', 'Night dose before bed'),
(3, '2025-08-04 13:00:00', 'Midday dose, slept in'),
(3, '2025-08-05 08:30:00', NULL),
(3, '2025-08-05 21:15:00', 'Bedtime dose'),
(3, '2025-08-06 09:00:00', 'Feeling rested'),
(3, '2025-08-07 08:45:00', NULL),
(3, '2025-08-07 20:45:00', 'Evening routine'),
(3, '2025-08-08 14:30:00', 'Late dose, forgot morning'),
(3, '2025-08-09 09:15:00', NULL),
(3, '2025-08-09 21:00:00', 'Night meds'),
(3, '2025-08-10 10:30:00', 'Weekend morning'),
(3, '2025-08-11 09:45:00', NULL),
(3, '2025-08-11 20:30:00', 'Evening dose'),
(3, '2025-08-12 08:30:00', 'Back to schedule'),
(3, '2025-08-12 21:15:00', 'Bedtime'),
(3, '2025-08-13 09:00:00', NULL),
(3, '2025-08-14 08:45:00', 'Morning routine'),
(3, '2025-08-14 20:45:00', 'Night dose'),
(3, '2025-08-15 15:00:00', 'Afternoon, running late'),
(3, '2025-08-16 09:15:00', NULL),
(3, '2025-08-16 21:00:00', 'Before bed'),
(3, '2025-08-17 08:30:00', 'Good morning'),
(3, '2025-08-18 09:00:00', NULL),
(3, '2025-08-18 20:30:00', 'Evening meds'),
(3, '2025-08-19 18:45:00', 'Late afternoon dose'),
(3, '2025-08-20 08:45:00', NULL),
(3, '2025-08-20 21:15:00', 'Nighttime'),
(3, '2025-08-21 09:15:00', 'Feeling stable'),
(3, '2025-08-22 08:30:00', NULL),
(3, '2025-08-22 20:45:00', 'Evening dose'),
(3, '2025-08-23 16:30:00', 'Late dose today'),
(3, '2025-08-24 10:00:00', 'Weekend sleep-in'),
(3, '2025-08-24 21:00:00', 'Night meds'),
(3, '2025-08-25 09:30:00', NULL),
(3, '2025-08-26 08:45:00', 'Monday start'),
(3, '2025-08-26 20:30:00', 'Evening routine'),
(3, '2025-08-27 09:00:00', NULL),
(3, '2025-08-27 21:15:00', 'Bedtime dose'),
(3, '2025-08-28 08:30:00', 'Mood is good'),
(3, '2025-08-29 14:00:00', 'Afternoon dose'),
(3, '2025-08-29 20:45:00', 'Night dose'),
(3, '2025-08-30 09:15:00', NULL),
(3, '2025-08-31 08:45:00', 'End of month'),
(3, '2025-08-31 21:00:00', 'Final evening dose of August');

-- medication_log_item entries
-- User 1 takes Sertraline (medication_id=1, 50mg) and Lithium (medication_id=6, 300mg)
-- Logs 1-31 are User 1
INSERT INTO medication_log_item (medication_log_id, medication_id, dosage, note) VALUES
-- Day 1
(1, 1, '50mg', NULL),
(1, 6, '300mg', NULL),
-- Day 2
(2, 1, '50mg', NULL),
(2, 6, '300mg', NULL),
-- Day 3
(3, 1, '50mg', NULL),
(3, 6, '300mg', NULL),
-- Day 4
(4, 1, '50mg', NULL),
(4, 6, '300mg', NULL),
-- Day 5
(5, 1, '50mg', NULL),
(5, 6, '300mg', NULL),
-- Day 6
(6, 1, '50mg', NULL),
(6, 6, '300mg', NULL),
-- Day 7
(7, 1, '50mg', NULL),
(7, 6, '300mg', NULL),
-- Day 8
(8, 1, '50mg', NULL),
(8, 6, '300mg', NULL),
-- Day 9
(9, 1, '50mg', NULL),
(9, 6, '300mg', NULL),
-- Day 10
(10, 1, '50mg', NULL),
(10, 6, '300mg', NULL),
-- Day 11
(11, 1, '50mg', NULL),
(11, 6, '300mg', NULL),
-- Day 12
(12, 1, '50mg', NULL),
(12, 6, '300mg', NULL),
-- Day 13
(13, 1, '50mg', NULL),
(13, 6, '300mg', NULL),
-- Day 14
(14, 1, '50mg', NULL),
(14, 6, '300mg', NULL),
-- Day 15
(15, 1, '50mg', NULL),
(15, 6, '300mg', NULL),
-- Day 16
(16, 1, '50mg', NULL),
(16, 6, '300mg', NULL),
-- Day 17
(17, 1, '50mg', NULL),
(17, 6, '300mg', NULL),
-- Day 18
(18, 1, '50mg', NULL),
(18, 6, '300mg', NULL),
-- Day 19
(19, 1, '50mg', NULL),
(19, 6, '300mg', NULL),
-- Day 20
(20, 1, '50mg', NULL),
(20, 6, '300mg', NULL),
-- Day 21
(21, 1, '50mg', NULL),
(21, 6, '300mg', NULL),
-- Day 22
(22, 1, '50mg', NULL),
(22, 6, '300mg', NULL),
-- Day 23
(23, 1, '50mg', NULL),
(23, 6, '300mg', NULL),
-- Day 24
(24, 1, '50mg', NULL),
(24, 6, '300mg', NULL),
-- Day 25
(25, 1, '50mg', NULL),
(25, 6, '300mg', NULL),
-- Day 26
(26, 1, '50mg', NULL),
(26, 6, '300mg', NULL),
-- Day 27
(27, 1, '50mg', NULL),
(27, 6, '300mg', NULL),
-- Day 28
(28, 1, '50mg', NULL),
(28, 6, '300mg', NULL),
-- Day 29
(29, 1, '50mg', NULL),
(29, 6, '300mg', NULL),
-- Day 30
(30, 1, '50mg', NULL),
(30, 6, '300mg', NULL),
-- Day 31
(31, 1, '50mg', NULL),
(31, 6, '300mg', NULL);

-- User 2 takes Fluoxetine (medication_id=2, 20mg) and Lamotrigine (medication_id=7, 100mg)
-- Logs 32-62 are User 2
INSERT INTO medication_log_item (medication_log_id, medication_id, dosage, note) VALUES
-- Day 1
(32, 2, '20mg', NULL),
(32, 7, '100mg', NULL),
-- Day 2
(33, 2, '20mg', NULL),
(33, 7, '100mg', NULL),
-- Day 3
(34, 2, '20mg', NULL),
(34, 7, '100mg', NULL),
-- Day 4
(35, 2, '20mg', NULL),
(35, 7, '100mg', NULL),
-- Day 5
(36, 2, '20mg', NULL),
(36, 7, '100mg', NULL),
-- Day 6
(37, 2, '20mg', NULL),
(37, 7, '100mg', NULL),
-- Day 7
(38, 2, '20mg', NULL),
(38, 7, '100mg', NULL),
-- Day 8
(39, 2, '20mg', NULL),
(39, 7, '100mg', NULL),
-- Day 9
(40, 2, '20mg', NULL),
(40, 7, '100mg', NULL),
-- Day 10
(41, 2, '20mg', NULL),
(41, 7, '100mg', NULL),
-- Day 11
(42, 2, '20mg', NULL),
(42, 7, '100mg', NULL),
-- Day 12
(43, 2, '20mg', NULL),
(43, 7, '100mg', NULL),
-- Day 13
(44, 2, '20mg', NULL),
(44, 7, '100mg', NULL),
-- Day 14
(45, 2, '20mg', NULL),
(45, 7, '100mg', NULL),
-- Day 15
(46, 2, '20mg', NULL),
(46, 7, '100mg', NULL),
-- Day 16
(47, 2, '20mg', NULL),
(47, 7, '100mg', NULL),
-- Day 17
(48, 2, '20mg', NULL),
(48, 7, '100mg', NULL),
-- Day 18
(49, 2, '20mg', NULL),
(49, 7, '100mg', NULL),
-- Day 19
(50, 2, '20mg', NULL),
(50, 7, '100mg', NULL),
-- Day 20
(51, 2, '20mg', NULL),
(51, 7, '100mg', NULL),
-- Day 21
(52, 2, '20mg', NULL),
(52, 7, '100mg', NULL),
-- Day 22
(53, 2, '20mg', NULL),
(53, 7, '100mg', NULL),
-- Day 23
(54, 2, '20mg', NULL),
(54, 7, '100mg', NULL),
-- Day 24
(55, 2, '20mg', NULL),
(55, 7, '100mg', NULL),
-- Day 25
(56, 2, '20mg', NULL),
(56, 7, '100mg', NULL),
-- Day 26
(57, 2, '20mg', NULL),
(57, 7, '100mg', NULL),
-- Day 27
(58, 2, '20mg', NULL),
(58, 7, '100mg', NULL),
-- Day 28
(59, 2, '20mg', NULL),
(59, 7, '100mg', NULL),
-- Day 29
(60, 2, '20mg', NULL),
(60, 7, '100mg', NULL),
-- Day 30
(61, 2, '20mg', NULL),
(61, 7, '100mg', NULL),
-- Day 31
(62, 2, '20mg', NULL),
(62, 7, '100mg', NULL);

-- User 3 takes Escitalopram (medication_id=3, 10mg) and Quetiapine (medication_id=8, 200mg)
-- User 3 has some logs with both meds, and separate evening Quetiapine doses
-- Logs 63-110 are User 3
INSERT INTO medication_log_item (medication_log_id, medication_id, dosage, note) VALUES
-- Day 1 morning
(63, 3, '10mg', NULL),
-- Day 1 evening (Quetiapine only)
(64, 8, '200mg', 'Bedtime dose'),
-- Day 2 morning
(65, 3, '10mg', NULL),
-- Day 3 morning
(66, 3, '10mg', NULL),
-- Day 3 evening
(67, 8, '200mg', 'Bedtime dose'),
-- Day 4 midday
(68, 3, '10mg', NULL),
-- Day 5 morning
(69, 3, '10mg', NULL),
-- Day 5 evening
(70, 8, '200mg', 'Bedtime dose'),
-- Day 6 morning
(71, 3, '10mg', NULL),
-- Day 7 morning
(72, 3, '10mg', NULL),
-- Day 7 evening
(73, 8, '200mg', 'Bedtime dose'),
-- Day 8 afternoon
(74, 3, '10mg', NULL),
(74, 8, '200mg', 'Taking both now'),
-- Day 9 morning
(75, 3, '10mg', NULL),
-- Day 9 evening
(76, 8, '200mg', 'Bedtime dose'),
-- Day 10 morning
(77, 3, '10mg', NULL),
-- Day 11 morning
(78, 3, '10mg', NULL),
-- Day 11 evening
(79, 8, '200mg', 'Bedtime dose'),
-- Day 12 morning
(80, 3, '10mg', NULL),
-- Day 12 evening
(81, 8, '200mg', 'Bedtime dose'),
-- Day 13 morning
(82, 3, '10mg', NULL),
-- Day 14 morning
(83, 3, '10mg', NULL),
-- Day 14 evening
(84, 8, '200mg', 'Bedtime dose'),
-- Day 15 afternoon
(85, 3, '10mg', NULL),
-- Day 16 morning
(86, 3, '10mg', NULL),
-- Day 16 evening
(87, 8, '200mg', 'Bedtime dose'),
-- Day 17 morning
(88, 3, '10mg', NULL),
-- Day 18 morning
(89, 3, '10mg', NULL),
-- Day 18 evening
(90, 8, '200mg', 'Bedtime dose'),
-- Day 19 late afternoon
(91, 3, '10mg', NULL),
(91, 8, '200mg', 'Taking both'),
-- Day 20 morning
(92, 3, '10mg', NULL),
-- Day 20 evening
(93, 8, '200mg', 'Bedtime dose'),
-- Day 21 morning
(94, 3, '10mg', NULL),
-- Day 22 morning
(95, 3, '10mg', NULL),
-- Day 22 evening
(96, 8, '200mg', 'Bedtime dose'),
-- Day 23 late afternoon
(97, 3, '10mg', NULL),
-- Day 24 morning
(98, 3, '10mg', NULL),
-- Day 24 evening
(99, 8, '200mg', 'Bedtime dose'),
-- Day 25 morning
(100, 3, '10mg', NULL),
-- Day 26 morning
(101, 3, '10mg', NULL),
-- Day 26 evening
(102, 8, '200mg', 'Bedtime dose'),
-- Day 27 morning
(103, 3, '10mg', NULL),
-- Day 27 evening
(104, 8, '200mg', 'Bedtime dose'),
-- Day 28 morning
(105, 3, '10mg', NULL),
-- Day 29 afternoon
(106, 3, '10mg', NULL),
-- Day 29 evening
(107, 8, '200mg', 'Bedtime dose'),
-- Day 30 morning
(108, 3, '10mg', NULL),
-- Day 31 morning
(109, 3, '10mg', NULL),
-- Day 31 evening
(110, 8, '200mg', 'Final bedtime dose of month');

-- ========================================
-- USEFUL QUERIES
-- ========================================

-- View all mood tags with their categories:
-- SELECT mt.name as mood, mc.name as category 
-- FROM mood_tag mt 
-- JOIN mood_category mc ON mt.mood_category_id = mc.mood_category_id 
-- ORDER BY mc.name, mt.name;

-- Get mood distribution by category:
-- SELECT mc.name as category, COUNT(mlmt.mood_tag_id) as usage_count
-- FROM mood_category mc
-- LEFT JOIN mood_tag mt ON mc.mood_category_id = mt.mood_category_id
-- LEFT JOIN mood_log_mood_tag mlmt ON mt.mood_tag_id = mlmt.mood_tag_id
-- GROUP BY mc.mood_category_id, mc.name
-- ORDER BY usage_count DESC;

-- View user's mood log with tags:
-- SELECT 
--     ml.mood_log_id,
--     u.email,
--     ml.mood_rating,
--     ml.note,
--     GROUP_CONCAT(mt.name SEPARATOR ', ') as mood_tags,
--     ml.created_at
-- FROM mood_log ml
-- JOIN user u ON ml.user_id = u.user_id
-- LEFT JOIN mood_log_mood_tag mlmt ON ml.mood_log_id = mlmt.mood_log_id
-- LEFT JOIN mood_tag mt ON mlmt.mood_tag_id = mt.mood_tag_id
-- WHERE u.user_id = 1
-- GROUP BY ml.mood_log_id
-- ORDER BY ml.created_at DESC;

-- Medication adherence report:
-- SELECT 
--     u.email,
--     m.name as medication,
--     DATE(ml.taken_at) as date,
--     SUM(CASE WHEN ml.taken = 1 THEN 1 ELSE 0 END) as doses_taken,
--     COUNT(*) as total_doses,
--     ROUND(SUM(CASE WHEN ml.taken = 1 THEN 1 ELSE 0 END) / COUNT(*) * 100, 2) as adherence_rate
-- FROM medication_log ml
-- JOIN user u ON ml.user_id = u.user_id
-- JOIN medication m ON ml.medication_id = m.medication_id
-- GROUP BY u.user_id, m.medication_id, DATE(ml.taken_at)
-- ORDER BY date DESC;

-- Sleep quality summary:
-- SELECT 
--     u.email,
--     AVG(sl.hours_slept) as avg_hours,
--     sqt.name as sleep_quality,
--     COUNT(*) as nights
-- FROM sleep_log sl
-- JOIN user u ON sl.user_id = u.user_id
-- JOIN sleep_quality_tag sqt ON sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
-- GROUP BY u.user_id, sqt.sleep_quality_tag_id
-- ORDER BY u.email, nights DESC;

-- Mood trends over time:
-- SELECT 
--     DATE(ml.created_at) as date,
--     AVG(ml.mood_rating) as avg_mood,
--     COUNT(*) as entries
-- FROM mood_log ml
-- WHERE ml.user_id = 1
-- GROUP BY DATE(ml.created_at)
-- ORDER BY date;

-- Correlation between sleep and mood:
-- SELECT 
--     DATE(sl.sleep_date) as date,
--     sl.hours_slept,
--     sqt.name as sleep_quality,
--     ml.mood_rating,
--     ml.note
-- FROM sleep_log sl
-- JOIN user u ON sl.user_id = u.user_id
-- LEFT JOIN mood_log ml ON u.user_id = ml.user_id 
--     AND DATE(sl.sleep_date) = DATE(ml.created_at)
-- JOIN sleep_quality_tag sqt ON sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
-- WHERE u.user_id = 1
-- ORDER BY date DESC;

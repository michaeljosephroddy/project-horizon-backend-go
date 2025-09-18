CREATE DATABASE IF NOT EXISTS project_horizon;
USE project_horizon;

-- Optional: Clean slate (use only in dev)
DROP TABLE IF EXISTS journal_entry_mood_tag;
DROP TABLE IF EXISTS user_medication;
DROP TABLE IF EXISTS medication;
DROP TABLE IF EXISTS journal_entry;
DROP TABLE IF EXISTS mood_tag;
DROP TABLE IF EXISTS user;

CREATE TABLE IF NOT EXISTS user (
    user_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS medication (
    medication_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS mood_tag (
    mood_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_medication (
    user_medication_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    medication_id BIGINT UNSIGNED NOT NULL,
    CONSTRAINT fk_um_user_medication FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_um_medication FOREIGN KEY (medication_id) REFERENCES medication(medication_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS journal_entry (
    journal_entry_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    mood_rating INT NOT NULL CHECK (mood_rating >= 1 AND mood_rating <= 10),
    note TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_journal_user FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    INDEX idx_created_at (created_at)
);

CREATE TABLE IF NOT EXISTS journal_entry_mood_tag (
    journal_entry_mood_tag_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE PRIMARY KEY, 
    journal_entry_id BIGINT UNSIGNED NOT NULL,
    mood_tag_id BIGINT UNSIGNED NOT NULL,
    CONSTRAINT fk_jem_journal_entry FOREIGN KEY (journal_entry_id) REFERENCES journal_entry(journal_entry_id) ON DELETE CASCADE,
    CONSTRAINT fk_jem_mood_tag FOREIGN KEY (mood_tag_id) REFERENCES mood_tag(mood_tag_id) ON DELETE CASCADE,
    INDEX idx_mood_tag (mood_tag_id)
);

-- Reset auto increment counters
ALTER TABLE user AUTO_INCREMENT = 1;
ALTER TABLE mood_tag AUTO_INCREMENT = 1;
ALTER TABLE journal_entry AUTO_INCREMENT = 1;
ALTER TABLE medication AUTO_INCREMENT = 1;
ALTER TABLE user_medication AUTO_INCREMENT = 1;
ALTER TABLE journal_entry_mood_tag AUTO_INCREMENT = 1;

-- Create a user for database access
CREATE USER IF NOT EXISTS 'demouser'@'localhost' IDENTIFIED BY 'demopassword';
GRANT ALL PRIVILEGES ON project_horizon.* TO 'demouser'@'localhost';
FLUSH PRIVILEGES;

-- Mock data for August 2025
-- Insert users
INSERT INTO user (email, password_hash, created_at) VALUES
('alice@example.com', '$2y$10$example.hash.1', '2025-07-15 10:00:00'),
('bob@example.com', '$2y$10$example.hash.2', '2025-07-20 14:30:00'),
('carol@example.com', '$2y$10$example.hash.3', '2025-07-25 09:15:00');

-- Insert medications
INSERT INTO medication (name) VALUES
('Sertraline'),
('Fluoxetine'),
('Escitalopram'),
('Bupropion'),
('Venlafaxine');

-- Insert mood tags
INSERT INTO mood_tag (name) VALUES
('Happy'),
('Sad'),
('Anxious'),
('Calm'),
('Energetic'),
('Tired'),
('Focused'),
('Overwhelmed'),
('Content'),
('Irritated');

-- Link users to medications
INSERT INTO user_medication (user_id, medication_id) VALUES
(1, 1), -- Alice takes Sertraline
(1, 4), -- Alice takes Bupropion
(2, 2), -- Bob takes Fluoxetine
(3, 3), -- Carol takes Escitalopram
(3, 5); -- Carol takes Venlafaxine

-- Journal entries for August 2025 (31 days)
INSERT INTO journal_entry (user_id, mood_rating, note, created_at) VALUES
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

-- Mood tag associations for journal entries (complete associations for all 93 entries)
INSERT INTO journal_entry_mood_tag (journal_entry_id, mood_tag_id) VALUES
-- Day 1 entries (Aug 1)
(1, 1), (1, 4), -- Happy, Calm
(2, 9), -- Content  
(3, 3), (3, 4), -- Anxious, Calm

-- Day 2 entries (Aug 2)
(4, 1), (4, 5), -- Happy, Energetic
(5, 2), (5, 6), -- Sad, Tired
(6, 1), (6, 9), -- Happy, Content

-- Day 3 entries (Aug 3)
(7, 4), (7, 9), -- Calm, Content
(8, 4), (8, 1), -- Calm, Happy
(9, 6), (9, 9), -- Tired, Content

-- Day 4 entries (Aug 4)
(10, 1), (10, 5), -- Happy, Energetic
(11, 9), (11, 7), -- Content, Focused
(12, 1), (12, 5), -- Happy, Energetic

-- Day 5 entries (Aug 5)
(13, 2), (13, 6), -- Sad, Tired
(14, 7), (14, 9), -- Focused, Content
(15, 8), (15, 3), -- Overwhelmed, Anxious

-- Day 6 entries (Aug 6)
(16, 4), (16, 7), -- Calm, Focused
(17, 9), -- Content
(18, 4), (18, 7), -- Calm, Focused

-- Day 7 entries (Aug 7)
(19, 5), (19, 1), -- Energetic, Happy
(20, 1), (20, 9), -- Happy, Content
(21, 4), (21, 9), -- Calm, Content

-- Day 8 entries (Aug 8)
(22, 9), (22, 7), -- Content, Focused
(23, 1), (23, 5), -- Happy, Energetic
(24, 6), (24, 8), -- Tired, Overwhelmed

-- Day 9 entries (Aug 9)
(25, 1), (25, 5), -- Happy, Energetic
(26, 9), (26, 6), -- Content, Tired
(27, 7), (27, 9), -- Focused, Content

-- Day 10 entries (Aug 10)
(28, 4), (28, 9), -- Calm, Content
(29, 1), (29, 5), -- Happy, Energetic
(30, 4), (30, 6), -- Calm, Tired

-- Day 11 entries (Aug 11)
(31, 1), (31, 9), -- Happy, Content
(32, 1), (32, 4), -- Happy, Calm
(33, 9), (33, 7), -- Content, Focused

-- Day 12 entries (Aug 12)
(34, 2), (34, 6), -- Sad, Tired
(35, 7), (35, 9), -- Focused, Content
(36, 8), (36, 3), -- Overwhelmed, Anxious

-- Day 13 entries (Aug 13)
(37, 7), (37, 9), -- Focused, Content
(38, 7), (38, 4), -- Focused, Calm
(39, 4), (39, 7), -- Calm, Focused

-- Day 14 entries (Aug 14)
(40, 1), (40, 5), -- Happy, Energetic
(41, 1), (41, 7), -- Happy, Focused
(42, 4), (42, 9), -- Calm, Content

-- Day 15 entries (Aug 15)
(43, 7), (43, 9), -- Focused, Content
(44, 9), -- Content
(45, 1), (45, 5), -- Happy, Energetic

-- Day 16 entries (Aug 16)
(46, 1), (46, 5), -- Happy, Energetic
(47, 1), (47, 9), -- Happy, Content
(48, 6), (48, 9), -- Tired, Content

-- Day 17 entries (Aug 17)
(49, 1), (49, 4), -- Happy, Calm
(50, 1), (50, 5), -- Happy, Energetic
(51, 4), (51, 9), -- Calm, Content

-- Day 18 entries (Aug 18)
(52, 4), (52, 9), -- Calm, Content
(53, 4), (53, 9), -- Calm, Content
(54, 1), (54, 9), -- Happy, Content

-- Day 19 entries (Aug 19)
(55, 7), (55, 9), -- Focused, Content
(56, 9), -- Content
(57, 9), -- Content

-- Day 20 entries (Aug 20)
(58, 1), (58, 5), -- Happy, Energetic
(59, 1), (59, 7), -- Happy, Focused
(60, 7), (60, 9), -- Focused, Content

-- Day 21 entries (Aug 21)
(61, 1), (61, 5), -- Happy, Energetic
(62, 7), (62, 9), -- Focused, Content
(63, 1), (63, 7), -- Happy, Focused

-- Day 22 entries (Aug 22)
(64, 6), (64, 8), -- Tired, Overwhelmed
(65, 1), (65, 5), -- Happy, Energetic
(66, 7), (66, 9), -- Focused, Content

-- Day 23 entries (Aug 23)
(67, 1), (67, 5), -- Happy, Energetic
(68, 1), (68, 9), -- Happy, Content
(69, 1), (69, 7), -- Happy, Focused

-- Day 24 entries (Aug 24)
(70, 4), (70, 9), -- Calm, Content
(71, 1), (71, 5), -- Happy, Energetic
(72, 4), (72, 9), -- Calm, Content

-- Day 25 entries (Aug 25)
(73, 1), (73, 9), -- Happy, Content
(74, 4), (74, 9), -- Calm, Content
(75, 1), (75, 5), -- Happy, Energetic

-- Day 26 entries (Aug 26)
(76, 7), (76, 9), -- Focused, Content
(77, 9), -- Content
(78, 5), (78, 7), -- Energetic, Focused

-- Day 27 entries (Aug 27)
(79, 1), (79, 7), -- Happy, Focused
(80, 7), (80, 9), -- Focused, Content
(81, 1), (81, 7), -- Happy, Focused

-- Day 28 entries (Aug 28)
(82, 1), (82, 5), -- Happy, Energetic
(83, 1), (83, 7), -- Happy, Focused
(84, 7), (84, 9), -- Focused, Content

-- Day 29 entries (Aug 29)
(85, 7), (85, 9), -- Focused, Content
(86, 1), (86, 5), -- Happy, Energetic
(87, 1), (87, 7), -- Happy, Focused

-- Day 30 entries (Aug 30)
(88, 1), (88, 5), -- Happy, Energetic
(89, 1), (89, 9), -- Happy, Content
(90, 1), (90, 7), -- Happy, Focused

-- Day 31 entries (Aug 31)
(91, 1), (91, 9), -- Happy, Content
(92, 1), (92, 9), -- Happy, Content
(93, 1), (93, 9); -- Happy, Content

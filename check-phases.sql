-- Check if template_phases table exists and has data
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = 'template_phases';

-- Check the structure
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'template_phases';

-- Check if there's any data
SELECT * FROM template_phases;

-- Check for a specific template
SELECT * FROM template_phases WHERE template_id = '3453b3b3-5ebd-4b4b-8921-3dfe414caf78';

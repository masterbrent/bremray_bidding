-- Add current_phase_id to jobs table
ALTER TABLE jobs 
ADD COLUMN current_phase_id VARCHAR(36) REFERENCES template_phases(id);

-- Create index for phase lookups
CREATE INDEX idx_jobs_current_phase_id ON jobs(current_phase_id);

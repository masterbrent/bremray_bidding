-- Create template_phases table
CREATE TABLE IF NOT EXISTS template_phases (
    id VARCHAR(36) PRIMARY KEY,
    template_id VARCHAR(36) NOT NULL REFERENCES job_templates(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phase_order INTEGER NOT NULL,
    description TEXT,
    UNIQUE(template_id, phase_order)
);

-- Create index for template lookups
CREATE INDEX idx_template_phases_template_id ON template_phases(template_id);

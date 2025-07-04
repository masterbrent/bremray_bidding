-- Add nickname column to items table
ALTER TABLE items ADD COLUMN IF NOT EXISTS nickname VARCHAR(255);

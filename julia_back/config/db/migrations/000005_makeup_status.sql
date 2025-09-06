-- Create makeup_status enum type
CREATE TYPE public.makeup_status AS ENUM (
    'pending',
    'present',
    'absent',
    'postponed',
    'exempt'
);

-- Add status column to makeups table with default value 'pending'
ALTER TABLE public.makeups 
ADD COLUMN status public.makeup_status NOT NULL DEFAULT 'pending';

CREATE TABLE public.exam_periods (
    exam_period_id BIGSERIAL PRIMARY KEY,    
    class_id BIGINT NOT NULL,                
    name TEXT NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,                
    end_date DATE NOT NULL,                  
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,

    CONSTRAINT fk_exam_periods_class
        FOREIGN KEY (class_id)
        REFERENCES public.classes (class_id)
        ON DELETE CASCADE
);

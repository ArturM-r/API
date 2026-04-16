CREATE TABLE public.apis (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              user_id UUID NOT NULL,
                              title TEXT NOT NULL,
                              completed BOOLEAN NOT NULL DEFAULT false,
                              created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

                              CONSTRAINT fk_todos_user
                                  FOREIGN KEY (user_id)
                                      REFERENCES public.users(id)
                                      ON DELETE CASCADE
);

CREATE INDEX idx_todos_user_id ON public.apis(user_id);

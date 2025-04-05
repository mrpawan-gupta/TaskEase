CREATE TABLE IF NOT EXISTS tasks (
                                     id VARCHAR(36) PRIMARY KEY,
                                     title VARCHAR(255) NOT NULL,
                                     description TEXT,
                                     status VARCHAR(20) NOT NULL CHECK (status IN ('Pending', 'InProgress', 'Completed')),
                                     created_at TIMESTAMP NOT NULL,
                                     updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
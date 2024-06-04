CREATE TABLE IF NOT EXISTS orders (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    game_id bigint NOT NULL,
    order_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    status VARCHAR(100) NOT NULL DEFAULT 'pending',  
    version integer NOT NULL DEFAULT 1,
    FOREIGN KEY (user_id) REFERENCES user_info(id),
    FOREIGN KEY (game_id) REFERENCES games(id)
);

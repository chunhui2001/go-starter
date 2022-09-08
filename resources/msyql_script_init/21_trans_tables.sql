
-- products
CREATE TABLE if not exists t_products (
   product_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   product_name VARCHAR(50),
   retail_price DOUBLE
) ENGINE = InnoDB;


CREATE TABLE if not exists t_customers (
  customer_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(50),
  last_name VARCHAR(50)
) ENGINE = InnoDB;

CREATE TABLE if not exists t_sales_orders (
  order_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  customer_id BIGINT,
  order_date DATETIME
) ENGINE = InnoDB;

CREATE TABLE if not exists t_sales_products (
  ref_id  BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT,
  product_id BIGINT,
  qty DOUBLE
) ENGINE = InnoDB;


-- INSERT INTO t_products (product_name, retail_price) VALUES ('5G SMARTPHONE', '850');
-- INSERT INTO t_products (product_name, retail_price) VALUES ('FITNESS TRACKER', '94');
-- INSERT INTO t_products (product_name, retail_price) VALUES ('1 LTR MINERAL WATER', '1.8');

-- INSERT INTO t_customers (first_name, last_name) VALUES ('JOHN', 'DOE');
-- INSERT INTO t_customers (first_name, last_name) VALUES ('MARY', 'SMITH');
-- INSERT INTO t_customers (first_name, last_name) VALUES ('STEVE', 'JACOBS');

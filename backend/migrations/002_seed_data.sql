-- 种子数据插入
-- 版本: 002
-- 创建时间: 2024-01-01 00:00:00

-- 插入管理员用户
INSERT INTO users (username, email, password, phone, real_name, status, role, balance, created_at, updated_at) 
VALUES 
('admin', 'admin@cloudbp.com', '$2a$10$N.zmdr9k7uOLQxaYdkFTxOHGJQJJ8A.O2JhbXPbBGtgcSFXzQKjUG', '13800138000', '系统管理员', 1, 'admin', 0.00, NOW(), NOW()),
('testuser', 'test@cloudbp.com', '$2a$10$N.zmdr9k7uOLQxaYdkFTxOHGJQJJ8A.O2JhbXPbBGtgcSFXzQKjUG', '13800138001', '测试用户', 1, 'user', 100.00, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 插入更多腾讯云Lighthouse产品
INSERT INTO products (provider_id, name, code, type, region, zone, cpu, memory, storage, storage_type, bandwidth, traffic, os, price, original_price, status, description, features, created_at, updated_at)
SELECT 
    p.id,
    '轻量应用服务器 4核8G',
    'lighthouse-4c8g',
    'lighthouse',
    'ap-guangzhou',
    'ap-guangzhou-3',
    4,
    8,
    180,
    'SSD',
    8,
    500,
    'Ubuntu 20.04',
    216.00,
    216.00,
    1,
    '适合高负载应用和服务',
    '{"support_docker": true, "support_ssh": true, "backup": true}',
    NOW(),
    NOW()
FROM providers p WHERE p.code = 'tencent' AND NOT EXISTS (
    SELECT 1 FROM products WHERE provider_id = p.id AND code = 'lighthouse-4c8g'
);

INSERT INTO products (provider_id, name, code, type, region, zone, cpu, memory, storage, storage_type, bandwidth, traffic, os, price, original_price, status, description, features, created_at, updated_at)
SELECT 
    p.id,
    '轻量应用服务器 8核16G',
    'lighthouse-8c16g',
    'lighthouse',
    'ap-guangzhou',
    'ap-guangzhou-3',
    8,
    16,
    300,
    'SSD',
    12,
    1000,
    'Ubuntu 20.04',
    432.00,
    432.00,
    1,
    '适合大型应用和数据库',
    '{"support_docker": true, "support_ssh": true, "backup": true}',
    NOW(),
    NOW()
FROM providers p WHERE p.code = 'tencent' AND NOT EXISTS (
    SELECT 1 FROM products WHERE provider_id = p.id AND code = 'lighthouse-8c16g'
);

-- 插入更多地域的产品
INSERT INTO products (provider_id, name, code, type, region, zone, cpu, memory, storage, storage_type, bandwidth, traffic, os, price, original_price, status, description, features, created_at, updated_at)
SELECT 
    p.id,
    '轻量应用服务器 1核2G (北京)',
    'lighthouse-1c2g-beijing',
    'lighthouse',
    'ap-beijing',
    'ap-beijing-3',
    1,
    2,
    50,
    'SSD',
    3,
    100,
    'Ubuntu 20.04',
    24.00,
    24.00,
    1,
    '适合个人开发者和小型应用 (北京地域)',
    '{"support_docker": true, "support_ssh": true, "backup": true}',
    NOW(),
    NOW()
FROM providers p WHERE p.code = 'tencent' AND NOT EXISTS (
    SELECT 1 FROM products WHERE provider_id = p.id AND code = 'lighthouse-1c2g-beijing'
);

INSERT INTO products (provider_id, name, code, type, region, zone, cpu, memory, storage, storage_type, bandwidth, traffic, os, price, original_price, status, description, features, created_at, updated_at)
SELECT 
    p.id,
    '轻量应用服务器 1核2G (上海)',
    'lighthouse-1c2g-shanghai',
    'lighthouse',
    'ap-shanghai',
    'ap-shanghai-2',
    1,
    2,
    50,
    'SSD',
    3,
    100,
    'Ubuntu 20.04',
    24.00,
    24.00,
    1,
    '适合个人开发者和小型应用 (上海地域)',
    '{"support_docker": true, "support_ssh": true, "backup": true}',
    NOW(),
    NOW()
FROM providers p WHERE p.code = 'tencent' AND NOT EXISTS (
    SELECT 1 FROM products WHERE provider_id = p.id AND code = 'lighthouse-1c2g-shanghai'
);

-- 插入更多系统配置
INSERT INTO configs (key, value, type, group_name, title, description, sort, status, created_at, updated_at)
VALUES 
('smtp_host', 'smtp.example.com', 'string', 'email', 'SMTP服务器', 'SMTP服务器地址', 1, 1, NOW(), NOW()),
('smtp_port', '587', 'int', 'email', 'SMTP端口', 'SMTP服务器端口', 2, 1, NOW(), NOW()),
('smtp_username', '', 'string', 'email', 'SMTP用户名', 'SMTP认证用户名', 3, 1, NOW(), NOW()),
('smtp_password', '', 'string', 'email', 'SMTP密码', 'SMTP认证密码', 4, 1, NOW(), NOW()),
('payment_methods', '["balance", "wechat", "alipay"]', 'json', 'payment', '支付方式', '支持的支付方式', 1, 1, NOW(), NOW()),
('min_recharge_amount', '10.00', 'float', 'payment', '最小充值金额', '用户最小充值金额', 2, 1, NOW(), NOW()),
('max_recharge_amount', '10000.00', 'float', 'payment', '最大充值金额', '用户最大充值金额', 3, 1, NOW(), NOW())
ON CONFLICT (key) DO NOTHING;
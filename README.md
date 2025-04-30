### **模块 1: user (用户模块)**

- **internal/domain/user**:
  - User 结构体 : 定义核心属性 (ID, Username, Email, Password (哈希), WalletAddr (登录用), PayoutWalletAddr (收款用), UserType, CreatedAt, UpdatedAt)。可以包含简单业务方法如 SetPayoutAddress(addr string) error (带校验)。
  - UserRepository 接口: 定义 FindByID, FindByEmail, FindByWalletAddr, Create, Update, Delete 等方法。
- **internal/module/user/repository**:
  - UserModel GORM 结构体: 包含 gorm.Model 和与 User 对应的字段，带 GORM 标签和数据库约束 (uniqueIndex 等)。
  - gormUserRepository 结构体: 实现 UserRepository 接口。
  - 实现 domainToModel, modelToDomain 映射函数。
  - 实现接口方法，封装 GORM 查询 (db.First, db.Where, db.Create, db.Save, db.Delete)。处理 gorm.ErrRecordNotFound -> apierror.NewNotFoundError。
- **internal/module/user/service**:
  - UserService 接口: 定义 Register, LoginByPassword,  GetUserByID,  等用例方法。
  - userService 结构体: 实现接口。注入 UserRepository, PasswordHasher (密码处理接口), TokenGenerator (JWT/Session接口)。
  - Register: 检查邮箱/地址是否已存在, 哈希密码, 创建 User 领域对象, 调用 userRepo.Create, **发布 UserRegisteredEvent** (如果用事件)。
  - UpdatePayoutAddress: 获取用户, 校验地址, 调用 userRepo.Update。
- **internal/module/user/handler**:
  - HTTPHandler 结构体: 注入 UserService。
  - 实现 Gin Handler 方法: 绑定校验 DTO, 从 Context 获取用户 ID (认证中间件设置), 调用 userService 方法, 处理错误, 返回 JSON 响应 (使用 UserOutput DTO)。
- **internal/module/router.go**:
  - RegisterRoutes 函数: 定义 /auth/..., /users/me, /users/me/payout-address 等路由，应用必要的中间件 (如认证)。

------



### **模块 2: product (商品模块)**

- **internal/domain/product**:
  - Product 结构体: ID, Name, Description, Price (common.Money), Stock (库存), Status (ProductStatus 值对象/枚举), CategoryID, Images ([]Image 值对象), CreatedAt, UpdatedAt。
  - Category 结构体 (实体)。
  - Image 结构体 (值对象)。
  - ProductRepository 接口: FindByID, Find, Create, Update, Delete, UpdateStock (可能需要原子操作)。
  - CategoryRepository 接口。
- **internal/module/product/repository**:
  - ProductModel, CategoryModel, ImageModel GORM 结构体。
  - gormProductRepository, gormCategoryRepository 实现。
  - 映射函数。
  - 使用 GORM 实现接口方法，注意 Preload 加载关联数据 (分类, 图片)。UpdateStock 可能需要使用 GORM 事务或数据库原子操作。
- **internal/module/product/service**:
  - ProductService 接口: CreateProduct, UpdateProduct, GetProductByID, ListProducts (带过滤/分页), ChangeProductStatus。
  - productService 结构体: 注入 ProductRepository, CategoryRepository。
  - 实现接口方法，调用 Repository。ListProducts 处理分页和过滤参数。
- **internal/module/product/handler**:
  - HTTPHandler: 注入 ProductService。
  - 实现 Handler: 调用 Service，返回 JSON。区分公共查询接口和管理后台接口 (后者需权限)。
- **internal/module/router.go**:
  - RegisterRoutes: 定义公共 GET /products, GET /products/{id} 和管理后台 POST /admin/products, PUT /admin/products/{id} 等路由。

------



### **模块 3: order (订单模块)**

- **internal/domain/order**:
  - Order 结构体 : ID, UserID, OrderSN (订单号), Status (OrderStatus 值对象/枚举), Items ([]OrderItem 实体), TotalPrice (common.Money), ShippingAddress (Address 值对象), CreatedAt, UpdatedAt。包含业务方法如 MarkAsPaid(), MarkAsShipped(), CalculateTotal()。
  - OrderItem 结构体 (实体): ID, OrderID, ProductID, ProductName, Price, Quantity。
  - OrderStatus 值对象/枚举。
  - OrderRepository 接口: FindByID, FindByUser, Create, Update。
- **internal/module/order/repository**:
  - OrderModel, OrderItemModel GORM 结构体。
  - gormOrderRepository 实现。
  - 映射函数。注意 Order 与 OrderItem 的一对多关系处理 (GORM Preload)。
- **internal/module/order/service**:
  - OrderService 接口: CreateOrder, GetOrderByID, ListUserOrders, MarkOrderAsPaid, MarkOrderAsShipped。
  - orderService 结构体: 注入 OrderRepository, ProductRepository (或 ProductService 用于查库存/价格), UserRepository (查地址), eventbus.Dispatcher。
  - CreateOrder: **开启事务**。校验用户、地址。检查商品库存和价格 (调用 product 模块)。创建 Order 和 OrderItem 领域对象。调用 orderRepo.Create 保存。**提交/回滚事务**。发布 OrderCreatedEvent。
  - MarkOrderAsPaid/Shipped: 加载 Order, 调用领域方法更新状态, orderRepo.Update 保存。通常由 payment 或 admin 模块调用。
- **internal/module/order/handler**:
  - HTTPHandler: 注入 OrderService。
  - 实现 Handler: 获取 UserID，调用 Service，返回 JSON。
- **internal/module/order/router.go**:
  - RegisterRoutes: 定义 POST /orders, GET /orders, GET /orders/{id}。可能需要管理后台查看所有订单的路由。

------



### **模块 4: payment (支付模块)**

- **internal/domain/payment**:
  - Transaction 结构体 : ID, UserID, OrderID (可选), Type (TxType 值对象/枚举: RECHARGE, PAYMENT, WITHDRAWAL, COMMISSION), Amount (common.Money), Currency, Status (TxStatus), TxHash (链上交易哈希), FromAddress, ToAddress, CreatedAt。
  - PaymentRepository 接口: CreateTransaction, FindTransactionByTxHash, FindPendingTransactions, UpdateTransactionStatus, FindOrderByAmount (用于匹配支付)。
- **internal/module/payment/repository**:
  - TransactionModel GORM 结构体。
  - gormPaymentRepository 实现。
  - 映射函数。FindOrderByAmount 需要特殊查询逻辑。
- **internal/module/payment/service**:
  - PaymentService 接口: GetPaymentDetails, ConfirmCryptoPayment, RequestWithdrawal, ProcessWithdrawal, GetUserBalance。
  - paymentService 结构体: 注入 PaymentRepository, OrderService, ReferralService (可选), UserRepository (获取 Payout 地址), Web3Client (平台发币用), config.PaymentConfig。
  - GetPaymentDetails: 确定收款地址 (配置), 计算需支付金额 (汇率), 返回给 OrderService。
  - ConfirmCryptoPayment (**由 Listener 调用**): **开启事务**。调用 paymentRepo.FindOrderByAmount 尝试匹配订单。创建 Transaction 记录。调用 orderService.MarkOrderAsPaid。调用 referralService.ProcessPurchase。**提交/回滚事务**。
  - RequestWithdrawal: 检查余额, **开启事务**。扣除余额 (WalletBalance), 创建 Pending 状态的 Transaction。**提交/回滚事务**。
  - ProcessWithdrawal (可能由 Admin 触发或定时任务): 获取 Pending 提现记录, **调用 Web3 接口转账到用户 Payout 地址**, 更新 Transaction 状态和 TxHash。
- **internal/module/payment/handler**:
  - HTTPHandler: 注入 PaymentService。
  - 实现 Handler: 获取 UserID，调用 Service，返回 JSON。
- **internal/module/payment/router.go**:
  - RegisterRoutes: 定义 POST /payments/withdraw, GET /payments/transactions, GET /wallet/balance。管理后台提现处理接口。

------



### **模块 5: referral (邀请/推荐模块)**

- **internal/domain/referral**:
  - ReferralProfile 结构体: 复杂状态，包含邀请码、上级、统计、余额、资格 (MaxCommissionDepth, CanWithdrawUSDT) 等。
  - RewardLedgerEntry 结构体 (实体): 记录奖励明细。
  - ReferralProfileRepository 接口: FindByID, FindByCode, Save, **GetAncestors** (核心！)。
  - RewardLedgerRepository 接口: Save, FindByUser。
  - RuleRepository 接口: 定义获取规则的方法。
  - MLMCommissionCalculator 领域服务: 封装多级佣金计算逻辑。
- **internal/module/referral/repository**:
  - ReferralProfileModel, RewardLedgerEntryModel GORM 结构体。
  - gormReferralProfileRepository 实现。**GetAncestors 需要特殊实现** (如递归 CTE)。
  - gormRewardLedgerRepository 实现。
  - configRuleRepository 或 dbRuleRepository 实现 (从配置或 DB 加载规则)。
  - 映射函数。
- **internal/module/referral/service**:
  - ReferralService 接口: EstablishReferral, ProcessPurchase, GetUserStats, UpdateRules (Admin)。
  - referralService 结构体: 注入 Repositories, RuleRepository, MLMCommissionCalculator。
  - EstablishReferral: 创建 Profile, 关联上级, 更新计数, 检查里程碑奖励。
  - ProcessPurchase (**由 PaymentService 调用**): **开启事务**。调用 MLMCommissionCalculator。保存 RewardLedgerEntry。更新 ReferralProfile 余额、累计消费、佣金深度等。**提交/回滚事务**。
- **internal/module/referral/handler**:
  - HTTPHandler: 注入 ReferralService。
  - 实现 Handler: 获取 UserID，调用 Service，返回 JSON。
- **internal/module/referral/router.go**:
  - RegisterRoutes: 定义 /referrals/my-stats, /referrals/my-code。管理后台规则设置接口。

------



### **模块 6: web3listener (区块链监听器)**

- **(非 HTTP 模块)**
- **listener.go**:
  - 依赖 ethclient.Client, PaymentService。
  - 初始化连接。
  - 设置 FilterQuery 监听 ERC20 Transfer 事件到平台地址。
  - 主循环 FilterLogs。
  - 管理 startBlock (持久化到文件或 DB)。
  - 实现**区块确认逻辑** (延迟处理或状态标记)。
  - 处理重连和错误。
- **handler.go**:
  - 解析 types.Log。
  - 在确认后，调用 paymentService.ConfirmCryptoPayment。
- **启动:** 在 cmd/app/main.go 中作为 goroutine 启动。

------



### **模块 7: notification (通知模块)**

- **internal/domain/notification**: (通常较简单)
  - (可选) NotificationLog 。
- **internal/module/notification/repository**:
  - (可选) gormNotificationLogRepository。
- **internal/module/notification/service**:
  - NotificationService 接口: SendWelcomeEmail, SendOrderPaidEmail, SendWithdrawalSuccessEmail。
  - emailNotificationService 实现: 注入邮件配置 (config.EmailConfig), 使用 gomail 等库发送邮件。
- **internal/module/notification/listener**: (如果使用事件驱动)
  - UserEventListener, OrderEventListener 等实现 eventbus.Listener。
  - 注入 NotificationService。
  - SubscribedTo() 返回监听的事件名。
  - Handle() 方法中调用 NotificationService 对应方法。
- **启动/注入:** 在 cmd/app/main.go 中创建 Service 实例，如果用事件则创建 Listener 并注册到 EventBus。

------



### **模块 8: admin (管理后台模块)**

- **核心职责:** 聚合各模块的管理功能，提供统一入口，处理管理员权限。
- **实现:**
  - 定义 /admin 前缀的路由。
  - Handler 层注入各个模块的 Service。
  - Handler 方法首先进行严格的**管理员权限校验** (通过中间件)。
  - 然后调用相应模块的 Service 方法执行操作 (如 userService.ListUsers, productService.CreateProduct, referralService.UpdateRules)。
  - 返回管理后台所需的特定格式数据。


### **支付，订单实现讨论结果**
- **相关细节：**
  - 用户支付的时候，会弹出相应的地址（也就是对应的收款地址）。
  - 在创建订单的时候，我们需要用户提供对应的汇款地址。
  - 也因此需要在订单的对应的表中需要添加对应的地址字段，时间戳等
  - 在支付模块中，我们通过监听对应的地址判断数额是否正确,在返回结果
  - 订单模块中调用对应支付模块的方法，得到对应的结果，再更新对应的状态
  
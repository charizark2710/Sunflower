[] Phải trả lại bộ nhớ nếu sử dụng với biến toàn cục (golang có garbage collection nhưng nếu phần bộ nhớ vẫn còn được sử dụng thì nó sẽ không tự xóa)

[] Khi sử dụng goroutine phải đặt size cho từng channel hoặc phải đảm bảo channel được đóng hoặc không còn data trong channel (memory leek)

[] Function luôn phải trả về error nếu có kết nối đến database hoặc rabbitmq hoặc với 1 bên thứ 3 nào khác.

[] Phải check error nếu function có trả về kiểu error

[] Interface model cho json và cho gorm object phải tách biệt, không được sử dụng chung

[] Nếu convert từ empty interface (interface{}) sang 1 kiểu bất kỳ nào đó phải check error.

[] Khi tạo API mới luôn defined trong RDIPs-BE/constant/ServiceConst  theo format 
<RequestMethod>.<Path> : serviceFunc

[] serviceFunc luôn theo dạng func(c *gin.Context) (commonModel.ResponseTemplate, error)

[] Tạo file service trong package services

[] Sử dụng package handler để chạy lệnh create và update cho database.

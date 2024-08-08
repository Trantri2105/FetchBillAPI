# Fetch camera bills API

## 1. Giới thiệu

Project cung cấp HTTP endpoints để trả về thông tin hóa đơn cho camera trong một khoảng thời gian nhất định

## 2.Setup
- Để sử dụng gmail để gửi file CSV, ta cần phải kích hoạt xác thực 2 bước và [tạo mật khẩu ứng dụng](https://support.google.com/mail/answer/185833?hl=vi)
- Khởi tạo các biến trong .env file
  - `FROM`: Mail gửi
  - `PASSWORD`: Mật khẩu ứng dụng
  - `TO`: Mail nhận
  - `SMTP_HOST`: Server smtp. VD `smtp.gmail.com` cho mail có đuôi `@gmail.com`
  - `SMTP_PORT`: Cổng smtp

## 3. Code chạy như nào

1. Query database để lấy danh sách record

```sh
select transaction_id, purchase_date_time, camera_sn, package_type
  from cam_bills
  where transaction_id is not null
  and payment_method = 'VIETTELPAY'
  and purchase_date_time > $1 AND purchase_date_time < $2
  and package_type in (select code from package_service where period > 3 and expired > 2595000)
```

2. Khởi động 1 goroutine để save record ra file csv và gửi file này qua mail với tên theo cú pháp sau:

```sh
./csv/bills_20241007-120000_20241015-120000.csv
```

## 4. Chi tiết api

### Lấy thông tin hóa đơn

- cURL :
  ```sh
  curl --location 'http://localhost:8080/' \
  --header 'Content-Type: application/json' \
  --data '{
      "start" : "30-10-2023 00:00:00",
      "end" : "29-07-2024 00:00:00",
  }'
  ```
- Body:

  - `start`: Ngày và giờ bắt đầu theo format `dd-MM-yyyy hh:mm:ss`. Ví dụ: `30-10-2023 00:00:00`
  - `end` : Ngày và giờ kết thúc theo format `dd-MM-yyyy hh:mm:ss`. Example: `29-07-2024 00:00:00`

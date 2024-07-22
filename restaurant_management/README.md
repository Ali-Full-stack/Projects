# Uyga vazifa: RabbitMQ Fanout Bilan E-commerce Buyurtma Boshqaruv Tizimi

## Maqsad
Ushbu vazifaning maqsadi `RabbitMQ`, `fanout` almashinuvlari bilan ishlash va buyurtmalarni boshqarish uchun sodda `REST API` yaratish.

## Talablar
1. **Buyurtmalarni boshqarish uchun REST API (Producer) yaratish**: 
    - Buyurtmalarni yaratish va olish uchun endpointlarga ega `REST API` ni amalga oshiring.
    - Yangi buyurtma yaratilganda, buyurtma tafsilotlarini `RabbitMQ` fanout` yuboring.
    - Har xil buyurtma holatlarini (statuslarini) qo'llab-quvvatlang (masalan, ko'rib chiqish, bajarilgan, bekor qilingan).

2. **Buyurtmalarni qayta ishlash uchun Worker (Consumer) yaratish:**: 
    - `Fanout` almashinuviga bog'langan `RabbitMQ` navbatini tinglaydigan consumer ni amalga oshiring.
    - Consumer buyurtmani qayta ishlashi kerak, shu jumladan:
        - Buyurtma tafsilotlarini `MongoDB` ga saqlash.
        - Buyurtma holatini tegishli ravishda yangilash.

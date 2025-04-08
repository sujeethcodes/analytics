
---

## 💡 Design Patterns

### 🔁 Singleton Pattern
We use the **Singleton Pattern** for the MongoDB connection (`MongoCon`) to ensure that:
- A single instance of the database connection is shared across the entire application.
- Resources are optimized, and connection pooling is used efficiently.
- All use cases and controllers access MongoDB through a shared instance injected via a container or constructor.

---

## 🧪 API Endpoints

### 1. Upload CSV File
**POST** `/upload-csv`  
Uploads and parses a CSV file containing sales data.

**Form Data:**
- `file` (CSV File)

**Example using `curl`:**
```bash
curl -F "file=@data.csv" http://localhost:8080/upload-csv

2. Get Revenue Analytics
GET /get-revenue?start_date=2024-01-01&end_date=2024-12-31&type=product
Returns total revenue grouped by type (optional):

product

category

region

empty for total revenue

3. Get Customer and Order Statistics
GET /get-customer-analysis?start_date=2024-01-01&end_date=2024-12-31
Returns:

Total Revenue

Total Orders

Unique Customers

Average Order Value

4. Get Profit Margin by Product
GET /get-other-calculations?start_date=2024-01-01&end_date=2024-12-31
Returns profit margin grouped by product:

Total Revenue

Total Cost

Profit Margin (percentage)

5. Refresh Analytics Data
GET /refresh
Refreshes the data from a fixed CSV path (you can configure the path in code).
⚠️ This endpoint deletes existing data before inserting new data.

📦 Features
Upload and parse CSV sales data.

Store and aggregate data using MongoDB.

Perform advanced analytics (grouped revenue, customer behavior, profit margin).

Use of context for timeout/cancellation handling.

🛠️ Technologies Used
Golang

Echo - HTTP Framework

MongoDB - NoSQL database for fast aggregation

CSV - Input data format

BSON Aggregation - For efficient analytics

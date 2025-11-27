# E-Commerce Fashion Store

A full-stack e-commerce platform for fashion products built with Go, React, and PostgreSQL.

## 🚀 Features

### Customer Features
- Browse products with advanced filtering
- Shopping cart and wishlist
- User authentication
- Checkout flow and order management

### Admin Features  
- Product and inventory management
- Order processing
- Analytics dashboard

## 🛠️ Tech Stack
- **Backend:** Go + PostgreSQL
- **Frontend:** React + Tailwind CSS
- **Admin:** React + Material-UI
- **DevOps:** Docker + docker-compose

## 📁 Project Structure

\`\`\`
ecommerce-fashion/
├── backend/           # Go API
├── frontend/          # Customer app
├── admin/             # Admin dashboard
└── docker-compose.yml
\`\`\`

## 🚦 Quick Start

\`\`\`bash
# Start all services
docker-compose up -d

# Run migrations  
docker-compose exec backend migrate -path migrations -database "$DATABASE_URL" up

# Access applications
# Customer: http://localhost:3000
# Admin: http://localhost:3001
# API: http://localhost:8080
\`\`\`

See individual README files in backend/, frontend/, and admin/ for detailed documentation.

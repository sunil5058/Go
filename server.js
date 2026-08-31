// server.js

require("dotenv").config();

const express = require("express");
const mongoose = require("mongoose");

const authRoutes = require("./routes/auth");

const app = express();

app.use(express.json());

mongoose
  .connect(process.env.MONGO_URI)
  .then(() => console.log("MongoDB connected"))
  .catch((error) => console.log(error));

app.use("/api/auth", authRoutes);

app.get("/", (req, res) => {
  res.json({
    message: "API is running",
  });
});

app.listen(5000, () => {
  console.log("Server running on port 5000");
});
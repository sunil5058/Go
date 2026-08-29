const express = require("express");
const cors = require("cors");
const helmet = require("helmet");
const morgan = require("morgan");

const routes = require("./routes");
const notFound = require("./middleware/notFound.middleware");
const errorHandler = require("./middleware/error.middleware");

const app = express();

// Security
app.use(helmet());
app.use(cors());

// Body parser
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Logger
app.use(morgan("dev"));

// Health check
app.get("/health", (req, res) => {
    res.status(200).json({
        success: true,
        message: "Server is running",
    });
});

// API routes
app.use("/api/v1", routes);

// 404
app.use(notFound);

// Error handler
app.use(errorHandler);

module.exports = app;
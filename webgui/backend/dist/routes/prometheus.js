"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = require("express");
const axios_1 = __importDefault(require("axios"));
const router = (0, express_1.Router)();
const PROMETHEUS_URL = process.env.PROMETHEUS_URL || 'http://prometheus-operated:9090/api/v1';
// Query Prometheus metrics
router.get('/', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const query = req.query.query;
        if (!query) {
            res.status(400).json({ message: 'Missing Prometheus query parameter' });
        }
        else {
            const url = `${PROMETHEUS_URL}/query`;
            const response = yield axios_1.default.get(url, {
                params: { query },
            });
            res.json(response.data);
        }
    }
    catch (error) {
        res.status(500).json({ message: 'Error fetching Prometheus metrics', error });
    }
}));
exports.default = router;

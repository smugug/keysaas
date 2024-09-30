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
const server_1 = require("../server");
const router = (0, express_1.Router)();
const BASE_URL = process.env.BASE_URL || 'https://kubernetes.default.svc/apis/keysaascontroller.keysaas/v1/namespaces/customer2/keysaases';
// Get all KeySaaS instances
router.get('/', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const response = yield axios_1.default.get(BASE_URL, {
            headers: { Authorization: `Bearer ${(0, server_1.getToken)()}` },
        });
        res.json(response.data);
    }
    catch (error) {
        res.status(500).json({ message: 'Error fetching KeySaaS instances', error });
    }
}));
// Create a new KeySaaS instance
router.post('/', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const response = yield axios_1.default.post(BASE_URL, req.body, {
            headers: { Authorization: `Bearer ${(0, server_1.getToken)()}` },
        });
        res.json(response.data);
    }
    catch (error) {
        res.status(500).json({ message: 'Error creating KeySaaS instance', error });
    }
}));
// Update a KeySaaS instance
router.put('/:name', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const url = `${BASE_URL}/${req.params.name}`;
        const response = yield axios_1.default.put(url, req.body, {
            headers: { Authorization: `Bearer ${(0, server_1.getToken)()}` },
        });
        res.json(response.data);
    }
    catch (error) {
        res.status(500).json({ message: `Error updating KeySaaS instance ${req.params.name}`, error });
    }
}));
// Delete a KeySaaS instance
router.delete('/:name', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const url = `${BASE_URL}/${req.params.name}`;
        const response = yield axios_1.default.delete(url, {
            headers: { Authorization: `Bearer ${(0, server_1.getToken)()}` },
        });
        res.json(response.data);
    }
    catch (error) {
        res.status(500).json({ message: `Error deleting KeySaaS instance ${req.params.name}`, error });
    }
}));
exports.default = router;

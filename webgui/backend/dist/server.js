"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getToken = void 0;
const express_1 = __importDefault(require("express"));
const keysaas_1 = __importDefault(require("./routes/keysaas"));
const prometheus_1 = __importDefault(require("./routes/prometheus"));
const fs_1 = __importDefault(require("fs"));
const app = (0, express_1.default)();
const port = process.env.PORT || 4000;
const token_path = process.env.TOKEN_PATH || '/var/run/secrets/kubernetes.io/serviceaccount/token';
app.use(express_1.default.json());
// Routes
app.use('/api/keysaas', keysaas_1.default);
app.use('/api/prometheus', prometheus_1.default);
app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});
const getToken = () => {
    try {
        const token = fs_1.default.readFileSync(token_path, 'utf8');
        return token;
    }
    catch (error) {
        console.error('Error reading service account token:', error);
        return '';
    }
};
exports.getToken = getToken;

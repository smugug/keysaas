import { Router } from 'express';
import { listKeysaasResources, createKeysaasResource, deleteKeysaasResource } from "./controllers";

const router = Router();

router.get('/keysaases', listKeysaasResources);


router.post('/keysaases', createKeysaasResource);


router.delete('/keysaases/:name', deleteKeysaasResource);

export { router as routes };

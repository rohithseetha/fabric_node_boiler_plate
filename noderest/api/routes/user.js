const express = require('express');

const router = express.Router();

//const mongoose = require('mongoose');

//const Contract = require('../model/contracts');
const read_file = require('../controllers/usercontroller');

// router.get('/',read_file.get_all_data);

router.get('/:ID',read_file.get_data);

router.post('/', read_file.load_data);


module.exports = router;

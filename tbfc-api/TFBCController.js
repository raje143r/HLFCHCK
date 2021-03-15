var express = require("express");
var router = express.Router();
var bodyParser = require("body-parser");

router.use(bodyParser.urlencoded({ extended: true }));
router.use(bodyParser.json());

const txSubmit = require("./invoke");
const txFetch = require("./query");

//var TFBC = require("./FabricHelper");

// Request LC
router.post("/requestLC", async function(req, res) {
  try {
    //req.body.expiryDate="21-03-2021";
    req.body.amount=parseInt(req.body.amount);
    req.body.paidamount=parseInt(req.body.paidamount);
    let result = await txSubmit.invoke("requestLC", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

// Issue LC
router.post("/issueLC", async function(req, res) {
  try {
    let result = await txSubmit.invoke("issueLC", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

// Accept LC
router.post("/acceptLC", async function(req, res) {
  try {
    let result = await txSubmit.invoke("acceptLC", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

// Reject LC
router.post("/rejectLC", async function(req, res) {
  try {
    let result = await txSubmit.invoke("rejectLC", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

// Get LC
router.post("/getLC", async function(req, res) {
  //TFBC.getLC(req, res); req.body.lcId
  try {
    let result = await txFetch.query("getLC", req.body.lcId);
    res.send(JSON.parse(result));
  } catch (err) {
    res.status(500).send(err);
  }
});

// Get LC history
router.post("/getLCHistory", async function(req, res) {
  try {
    let result = await txFetch.query("getLCHistory", req.body.lcId);
    res.send(JSON.parse(result));
  } catch (err) {
    res.status(500).send(err);
  }
});

// Update Carrier
router.post("/updateCarrier", async function(req, res) {
  try {
    let result = await txSubmit.invoke("updateCarrier", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

// Ship Goods
router.post("/shipGoods", async function(req, res) {
  try {
    let result = await txSubmit.invoke("shipGoods", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

// Pay Seller
router.post("/paySeller", async function(req, res) {
  try {
    req.body.paidamount = parseInt(req.body.paidamount);
    let result = await txSubmit.invoke("paySeller", JSON.stringify(req.body));
    res.send(result);
  } catch (err) {
    res.status(500).send(err);
  }
});

module.exports = router;

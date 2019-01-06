package sticker

import (
	"fmt"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"math"
	"os/exec"
	"time"
)

const (
	LEFT_MOTOR  = 0
	RIGHT_MOTOR = 1

	MOTOR_HAT_ADDRESS = 0x60
	MOTOR_SPEED       = 100

	SERVO_HAT_ADDRESS                = 0x40
	SERVO_LOWER_ANGLE_BOUND          = -90.0
	SERVO_UPPER_ANGLE_BOUND          = 90.0
	SERVO_ANGLE_RANGE                = SERVO_UPPER_ANGLE_BOUND - SERVO_LOWER_ANGLE_BOUND
	SERVO_LOWER_ANGLE_PULSE_WIDTH_MS = 1.0
	SERVO_UPPER_ANGLE_PULSE_WIDTH_MS = 2.0
	SERVO_PULSE_WIDTH_RANGE          = SERVO_UPPER_ANGLE_PULSE_WIDTH_MS - SERVO_LOWER_ANGLE_PULSE_WIDTH_MS

	SERVO_INCREMENT = 1.0

	UPPER_PERIOD_MS          = SERVO_UPPER_ANGLE_PULSE_WIDTH_MS + 5
	SERVO_FREQUENCY_HZ       = 1.0 / (UPPER_PERIOD_MS / 1000.0)
	PERIOD_PER_COUNT_MS      = UPPER_PERIOD_MS / 4096
	LOWER_BOUND_COUNT        = SERVO_LOWER_ANGLE_PULSE_WIDTH_MS / PERIOD_PER_COUNT_MS
	ELEVATION_START_POSITION = 60.0
)

const (
	Azimuth   ServoType = 8
	Elevation ServoType = 4
)

const (
	Release ServoMode = iota
	Increase
	Decrease
)

type ServoType int
type ServoMode int

type Robot struct {
	motorServoDriver *i2c.AdafruitMotorHatDriver
	servoModes       map[ServoType]ServoMode
	servoAngle       map[ServoType]float64
}

func NewRobot() *Robot {
	r := raspi.NewAdaptor()
	robot := Robot{
		i2c.NewAdafruitMotorHatDriver(r),
		map[ServoType]ServoMode{
			Azimuth:   Release,
			Elevation: Release,
		},
		map[ServoType]float64{
			Azimuth:   0.0,
			Elevation: ELEVATION_START_POSITION,
		},
	}

	robot.motorServoDriver.SetMotorHatAddress(MOTOR_HAT_ADDRESS)
	robot.motorServoDriver.SetServoHatAddress(SERVO_HAT_ADDRESS)
	return &robot
}

func (r Robot) Start() {
	r.motorServoDriver.Start()
	r.motorServoDriver.SetDCMotorSpeed(LEFT_MOTOR, MOTOR_SPEED)
	r.motorServoDriver.SetDCMotorSpeed(RIGHT_MOTOR, MOTOR_SPEED)
	r.motorServoDriver.SetServoMotorFreq(SERVO_FREQUENCY_HZ)
	go updateServo(&r)
	startImageCapture()
}

func (r Robot) turnLeft() {
	r.motorServoDriver.RunDCMotor(LEFT_MOTOR, i2c.AdafruitBackward)
	r.motorServoDriver.RunDCMotor(RIGHT_MOTOR, i2c.AdafruitForward)
}

func (r Robot) turnRight() {
	r.motorServoDriver.RunDCMotor(LEFT_MOTOR, i2c.AdafruitForward)
	r.motorServoDriver.RunDCMotor(RIGHT_MOTOR, i2c.AdafruitBackward)
}

func (r Robot) moveForward() {
	r.motorServoDriver.RunDCMotor(LEFT_MOTOR, i2c.AdafruitForward)
	r.motorServoDriver.RunDCMotor(RIGHT_MOTOR, i2c.AdafruitForward)
}

func (r Robot) moveBackward() {
	r.motorServoDriver.RunDCMotor(LEFT_MOTOR, i2c.AdafruitBackward)
	r.motorServoDriver.RunDCMotor(RIGHT_MOTOR, i2c.AdafruitBackward)
}

func (r Robot) stopMove() {
	r.motorServoDriver.RunDCMotor(LEFT_MOTOR, i2c.AdafruitRelease)
	r.motorServoDriver.RunDCMotor(RIGHT_MOTOR, i2c.AdafruitRelease)
}

func (r Robot) lookUp() {
	r.setServo(Elevation, Increase)
}

func (r Robot) lookDown() {
	r.setServo(Elevation, Decrease)
}

func (r Robot) lookLeft() {
	r.setServo(Azimuth, Increase)
}

func (r Robot) lookRight() {
	r.setServo(Azimuth, Decrease)
}

func (r Robot) stopLook() {
	r.setServo(Azimuth, Release)
	r.setServo(Elevation, Release)
}

func (r Robot) setServo(stype ServoType, smode ServoMode) {
	r.servoModes[stype] = smode
}

func angleToCounts(angle float64) int32 {
	angleRatio := (angle - SERVO_LOWER_ANGLE_BOUND) / SERVO_ANGLE_RANGE
	periodMs := SERVO_PULSE_WIDTH_RANGE*angleRatio + SERVO_LOWER_ANGLE_PULSE_WIDTH_MS
	counts := periodMs / PERIOD_PER_COUNT_MS
	return int32(math.Floor(counts))
}

func updateServoAngles(robot *Robot) {
	azimuthCounts := angleToCounts(robot.servoAngle[Azimuth])
	robot.motorServoDriver.SetServoMotorPulse(byte(Azimuth), 0, azimuthCounts)
	elevationCounts := angleToCounts(robot.servoAngle[Elevation])
	robot.motorServoDriver.SetServoMotorPulse(byte(Elevation), 0, elevationCounts)
}

func updateServo(robot *Robot) {
	for {
		switch robot.servoModes[Azimuth] {
		case Increase:
			increaseServo(robot, Azimuth)
		case Decrease:
			decreaseServo(robot, Azimuth)
		}

		switch robot.servoModes[Elevation] {
		case Increase:
			decreaseServo(robot, Elevation)
		case Decrease:
			increaseServo(robot, Elevation)
		}
		updateServoAngles(robot)
		time.Sleep(20 * time.Millisecond)
	}
}

func increaseServo(robot *Robot, sType ServoType) {
	robot.servoAngle[sType] = math.Min(
		robot.servoAngle[sType]+SERVO_INCREMENT, SERVO_UPPER_ANGLE_BOUND)
}

func decreaseServo(robot *Robot, sType ServoType) {
	robot.servoAngle[sType] = math.Max(
		robot.servoAngle[sType]-SERVO_INCREMENT, SERVO_LOWER_ANGLE_BOUND)
}

func startImageCapture() {
	cmd := exec.Command("/host/sticker/object_detection/continuous_detection.py")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error:  %v", err)
	}
}

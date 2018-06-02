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

	SERVO_LOWER_ANGLE_BOUND          = -90.0
	SERVO_UPPER_ANGLE_BOUND          = 90.0
	SERVO_ANGLE_RANGE                = SERVO_UPPER_ANGLE_BOUND - SERVO_LOWER_ANGLE_BOUND
	SERVO_LOWER_ANGLE_PULSE_WIDTH_MS = 1.0
	SERVO_UPPER_ANGLE_PULSE_WIDTH_MS = 2.0
	SERVO_PULSE_WIDTH_RANGE          = SERVO_UPPER_ANGLE_PULSE_WIDTH_MS - SERVO_LOWER_ANGLE_PULSE_WIDTH_MS

	SERVO_INCREMENT = 1.0

	UPPER_PERIOD_MS     = SERVO_UPPER_ANGLE_PULSE_WIDTH_MS + 5
	SERVO_FREQUENCY_HZ  = 1.0 / (UPPER_PERIOD_MS / 1000.0)
	PERIOD_PER_COUNT_MS = UPPER_PERIOD_MS / 4096
	LOWER_BOUND_COUNT   = SERVO_LOWER_ANGLE_PULSE_WIDTH_MS / PERIOD_PER_COUNT_MS
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
			Elevation: 0.0,
		},
	}

	robot.motorServoDriver.SetMotorHatAddress(0x60)
	robot.motorServoDriver.Start()
	robot.motorServoDriver.SetDCMotorSpeed(LEFT_MOTOR, 100)
	robot.motorServoDriver.SetDCMotorSpeed(RIGHT_MOTOR, 100)
	robot.motorServoDriver.SetServoHatAddress(0x40)
	robot.motorServoDriver.SetServoMotorFreq(SERVO_FREQUENCY_HZ)

	go update_servo(&robot)
	start_image_capture()

	return &robot
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
	r.set_servo(Elevation, Increase)
}

func (r Robot) lookDown() {
	r.set_servo(Elevation, Decrease)
}

func (r Robot) lookLeft() {
	r.set_servo(Azimuth, Increase)
}

func (r Robot) lookRight() {
	r.set_servo(Azimuth, Decrease)
}

func (r Robot) stopLook() {
	r.set_servo(Azimuth, Release)
	r.set_servo(Elevation, Release)
}

func (r Robot) set_servo(stype ServoType, smode ServoMode) {
	r.servoModes[stype] = smode
}

func angle_to_counts(angle float64) int32 {
	angleRatio := (angle - SERVO_LOWER_ANGLE_BOUND) / SERVO_ANGLE_RANGE
	period_ms := SERVO_PULSE_WIDTH_RANGE*angleRatio + SERVO_LOWER_ANGLE_PULSE_WIDTH_MS
	counts := period_ms / PERIOD_PER_COUNT_MS
	return int32(math.Floor(counts))
}

func update_servo_angles(robot *Robot) {
	azimuthCounts := angle_to_counts(robot.servoAngle[Azimuth])
	robot.motorServoDriver.SetServoMotorPulse(byte(Azimuth), 0, azimuthCounts)
	elevationCounts := angle_to_counts(robot.servoAngle[Elevation])
	robot.motorServoDriver.SetServoMotorPulse(byte(Elevation), 0, elevationCounts)
}

func update_servo(robot *Robot) {
	for {
		switch robot.servoModes[Azimuth] {
		case Increase:
			increase_servo(robot, Azimuth)
		case Decrease:
			decrease_servo(robot, Azimuth)
		}

		switch robot.servoModes[Elevation] {
		case Increase:
			decrease_servo(robot, Elevation)
		case Decrease:
			increase_servo(robot, Elevation)
		}
		update_servo_angles(robot)
		time.Sleep(20 * time.Millisecond)
	}
}

func increase_servo(robot *Robot, sType ServoType) {
	robot.servoAngle[sType] = math.Min(
		robot.servoAngle[sType]+SERVO_INCREMENT, SERVO_UPPER_ANGLE_BOUND)
}

func decrease_servo(robot *Robot, sType ServoType) {
	robot.servoAngle[sType] = math.Max(
		robot.servoAngle[sType]-SERVO_INCREMENT, SERVO_LOWER_ANGLE_BOUND)
}

func start_image_capture() {
	cmd := exec.Command("raspistill", "-o", "assets/html/image.jpg", "-w", "400", "-h", "300", "-tl", "500", "-t", "0")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error:  %v", err)
	}
}

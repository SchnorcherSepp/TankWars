#!/usr/bin/env python

"""
This file contains an example in Python for an AI controlled client.
Use this example to program your own AI in Python.
"""

import json
import socket
import time
from threading import Lock

# CONFIG
TCP_IP = '127.0.0.1'
TCP_PORT = 3333

# TCP connection
conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
conn.connect((TCP_IP, TCP_PORT))

# read first line
conn_make_file = conn.makefile()
hello = conn_make_file.readline()
print(hello)

# thread lock
lock = Lock()


# ------ Helper ------------------------------------------------------------------------------------------------------ #


def command(cmd):
    lock.acquire()  # <---- LOCK

    # remove protocol break
    cmd = cmd.replace('\n', '')
    cmd = cmd.replace('\r', '')

    # send command
    conn.send(bytes(cmd, 'utf8') + b'\n')
    print("SEND:", cmd)  # DEBUG !!!

    # read response
    resp = conn_make_file.readline()
    resp = resp.replace('\n', '')
    resp = resp.replace('\r', '')
    print("RESP:", resp)  # DEBUG !!!

    lock.release()  # <---- UNLOCK

    # return
    return resp


# ---------------- GETTER --------------------------------------------------------------------------------------------#


# my_name returns the active player of this connection (RedTank or BlueTank)
def my_name():
    return command("MyName")


# game_status returns a json with all world data.
def game_status():
    return command("GameStatus")


# tank_status returns a json with all data of a requested tank.
def tank_status(tank_id):
    return command("TankStatus %s" % tank_id)


# close_targets returns all objects in the world that are theoretical in weapon range.
# The weapon type is irrelevant (WeaponCannon or WeaponArtillery) and the angle of the tank is ignored.
# The list is sorted by distance (from the closest to the farthest).
def close_targets(tank_id, f1, f2, f3, f4, f5):
    return command("CloseTargets %s %s %s %s %s %s" % (tank_id, f1, f2, f3, f4, f5))


# possible_targets extends CloseTargets.
# It only returns objects that can actually be attacked,depending on the weapon type.
# However, it may be necessary for the battle tank to change its angle.
# The list is sorted by the rotation required to reach the target.
def possible_targets(tank_id, f1, f2, f3, f4, f5):
    return command("PossibleTargets %s %s %s %s %s %s" % (tank_id, f1, f2, f3, f4, f5))


# ---------------- SETTER --------------------------------------------------------------------------------------------#


# exit_server kills the server (for tests only).
def exit_server():
    return command("Exit")


# buy_tank buy a new tank and place it near the home base.
def buy_tank(armor, damage, weapon):
    return command("BuyTank %d %d %s" % (armor, damage, weapon))


# fire creates a new projectile.
# The attributes fireAngle and distance determine the direction and distance of the shot.
# Cannons can fire in vehicle angle only.
# The distance is limited by the weapon range.
def fire(tank_id, angle, distance):
    return command("Fire %s %d %d" % (tank_id, angle, distance))


# fire_at is a wrapper for Fire() and convert the position to fireAngle and distance.
def fire_at(tank_id, x, y):
    return command("FireAt %s %d %d" % (tank_id, x, y))


# forward send the tank forward.
def forward(tank_id):
    return command("Forward %s" % tank_id)


# backward send the tank back.
def backward(tank_id):
    return command("Backward %s" % tank_id)


# stop the movement.
# Weapons can only build up when the tank is stationary
def stop(tank_id):
    return command("Stop %s" % tank_id)


# left turn the tank direction 45° left.
def left(tank_id):
    return command("Left %s" % tank_id)


# right turn the tank direction 45° right.
def right(tank_id):
    return command("Right %s" % tank_id)


# set_macro_move_to sets a special macro with a position that is called with every update.
def set_macro_move_to(tank_id, x, y):
    return command("SetMacroMoveTo %s %d %d" % (tank_id, x, y))


# set_macro sets a macro that is called with every update.
def set_macro(tank_id, macro):
    return command("SetMacro %s %s" % (tank_id, macro))


# ----- MY AI ---------------------------------------------------------------------------------------- #


if __name__ == '__main__':

    # get my player
    me = my_name()
    print(">", me)

    # request world status
    json_str = game_status()
    world = json.loads(json_str)  # UPDATE WORLD !!!
    print(">", world)

    # find my tanks in 'world'
    cannon_id = "-1"
    rocket_id = "-1"

    for tank in world["tanks"]:
        if tank["owner"] == me:
            # I am the owner of this tank
            if tank["weapon"]["typ"] == "Tank":
                # found battle tank
                cannon_id = tank["id"]

            if tank["weapon"]["typ"] == "RocketLauncher":
                # found rocket launcher
                rocket_id = tank["id"]

    # print id of my tanks
    print(">", "cannon_id", cannon_id)
    print(">", "rocket_id", rocket_id)

    # get center position
    # (use world.ScreenWidth and world.ScreenHeight)
    center_x = world["screenWidth"] / 2
    center_y = world["screenHeight"] / 2
    if me == "red":
        center_x -= 150  # position if I am red
    else:
        center_x += 150  # position if I am blue

    # use MACRO: MoveTo to move tank to the center
    err = set_macro_move_to(cannon_id, center_x, center_y)
    print(">", "move tank", err)

    err = set_macro_move_to(rocket_id, center_x, center_y)
    print(">", "move rocket", err)

    # wait for tanks
    # (wait for iteration 800)
    while world["iteration"] < 800:
        world = json.loads(game_status())  # UPDATE WORLD !!!
        time.sleep(5)

    # use MACRO: GuardMode to defend the center
    err = set_macro(cannon_id, "GuardMode")
    print(">", "macro tank", err)

    err = set_macro(rocket_id, "GuardMode")
    print(">", "macro rocket", err)

    # wait for more money
    while world["cashRed"] < 101 and world["cashBlue"] < 101:
        world = json.loads(game_status())  # UPDATE WORLD !!!
        time.sleep(5)

    # buy Artillery
    err = buy_tank(5, 70, "Artillery")
    print(">", "buy", err)

    # find my tanks in 'world'
    artillery_id = "-1"

    world = json.loads(game_status())  # UPDATE WORLD !!!
    for tank in world["tanks"]:
        if tank["owner"] == me:
            # I am the owner of this tank
            if tank["weapon"]["typ"] == "Artillery":
                # found battle tank
                artillery_id = tank["id"]

    print(">", "artillery_id", artillery_id)

    # use MACRO: MacroAttackMove to attack with artillery
    err = set_macro(artillery_id, "AttackMove")
    print(">", "macro artillery", err)

package ch.unisg.ics.interactions;

import cartago.OPERATION;
import cartago.ObsProperty;
import org.hyperagents.yggdrasil.cartago.artifacts.HypermediaArtifact;
import org.hyperagents.yggdrasil.cartago.artifacts.HypermediaTDArtifact;


public class SecuritySystem extends HypermediaTDArtifact {


  public void init() {
    defineObsProperty("isLocked", true);
  }


  @OPERATION
  public void unlockRoom() {
    ObsProperty roomIsLocked = getObsProperty("isLocked");

    if (!roomIsLocked.booleanValue()) {
      log("Room is already unlocked");
      return;
    }

    roomIsLocked.updateValue(false);

  }

  @OPERATION
  public void lockRoom() {
    ObsProperty roomIsLocked = getObsProperty("isLocked");

    if (roomIsLocked.booleanValue()) {
      log("Room is already locked");
      return;
    }

    roomIsLocked.updateValue(true);

  }

  @Override
  protected void registerInteractionAffordances() {
    // Register one action affordance with an input schema
    registerActionAffordance("http://example.org/unlockRoom", "unlockRoom", "/unlockRoom");
    registerActionAffordance("http://example.org/lockRoom", "lockRoom", "/lockRoom");
  }

}

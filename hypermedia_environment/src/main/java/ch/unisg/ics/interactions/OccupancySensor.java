package ch.unisg.ics.interactions;

import cartago.OPERATION;
import cartago.ObsProperty;
import org.hyperagents.yggdrasil.cartago.artifacts.HypermediaArtifact;
import org.hyperagents.yggdrasil.cartago.artifacts.HypermediaTDArtifact;

public class OccupancySensor extends HypermediaTDArtifact {


  public void init() {
    defineObsProperty("fallDetected", false);
    defineObsProperty("roomIsOccupied", false);
  }


  @OPERATION
  public void triggerFallDetected() {
    ObsProperty fallDetected = getObsProperty("fallDetected");
    fallDetected.updateValue(true);
    log("FALL DETECTED");
  }

  @Override
  protected void registerInteractionAffordances() {
    // Register one action affordance with an input schema
    registerActionAffordance("http://example.org/triggerFallDetected", "triggerFallDetected", "/triggerFallDetected");
  }

}
